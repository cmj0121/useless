#include <stdio.h>
#include <stdlib.h>
#include <elf.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>

#define DEBUG(msg, ...)	\
	do {				\
		fprintf(stderr, "[%s #%04d] " msg "\n", __FILE__, __LINE__, ##__VA_ARGS__); \
	} while (0)

#define BYTE( expr ) ( (char)((expr) & 0xFF))

#define KEY	0x4A


int main(int argc, char* argv[]) {
	int fd = -1;
	char src[] = "a.out", dst[] = "b.out";

	if (0 > (fd = open(src, O_RDONLY))) {
		DEBUG("Cannot open file %s", argv[1]);
		return -1;
	}
	void *ptr = NULL;
	off_t total_filesz = lseek(fd, 0, SEEK_END);

	if (MAP_FAILED == (ptr = mmap(NULL, total_filesz, PROT_READ|PROT_WRITE, MAP_PRIVATE, fd, 0))) {
		DEBUG("Cannot map the file to memory");
		return -1;
	}

	Elf64_Ehdr *hdr = ptr;
	char decoder[] = {
		/* mov	rax	Entry      */
		0xb8, 0x00, 0x00, 0x00, 0x00,
		/* mov	rcx	filesz     */
		0xb9, 0x00, 0x00, 0x00, 0x00,
		/* mov	bl	KEY        */
		0xB3, 0x00,
		/* xor	byte [rax] KEY */
		//0x30, 0x18,
		/* inc	rax            */
		0xFF, 0xC0,
		/* loopz               */
		0XE1, 0XFB,
		/* jmp	to Entry       */
		0xE9, 0x00, 0x00, 0x00, 0x00,
	};

	for (int i = 0; i < hdr->e_phnum; i++) {
		long offset = i * sizeof(Elf64_Phdr) + hdr->e_phoff;
		Elf64_Phdr *phdr = ptr + offset;

		if (phdr->p_type == PT_LOAD &&  phdr->p_flags == (PF_X+PF_R)) {
			/* Find the necessary program header */
			uint64_t offset = phdr->p_offset;
			uint64_t filesz = phdr->p_filesz;
			long entry_offset = (long)0 - filesz - sizeof(decoder) + (hdr->e_entry - offset);

			/* Update the offset */
			decoder[1] = BYTE(offset);
			decoder[2] = BYTE(offset >> 8);
			decoder[3] = BYTE(offset >> 16);
			decoder[4] = BYTE(offset >> 24);
			/* Update the filesz */
			decoder[6] = BYTE(filesz);
			decoder[7] = BYTE(filesz >> 8);
			decoder[8] = BYTE(filesz >> 16);
			decoder[9] = BYTE(filesz >> 24);
			/* Update the entry */
			decoder[sizeof(decoder)-4] = BYTE(entry_offset);
			decoder[sizeof(decoder)-3] = BYTE(entry_offset >> 8);
			decoder[sizeof(decoder)-2] = BYTE(entry_offset >> 16);
			decoder[sizeof(decoder)-1] = BYTE(entry_offset >> 24);

			DEBUG("Program size : %X, offset : %X", filesz, offset);
			for (uint64_t off = 0; off < phdr->p_filesz; ++off) {
				//((char *)ptr)[offset + off] ^= KEY;
			}

			for (int i = 0; i < sizeof(decoder); ++i) {
				((char *)ptr)[phdr->p_filesz + phdr->p_offset + i] = decoder[i];
			}

			DEBUG("New entry : %X", phdr->p_filesz + phdr->p_offset);

			/* Set writable */
			phdr->p_flags = PF_R + PF_W + PF_X;
			/* Change the entry address */
			hdr->e_entry = phdr->p_filesz + phdr->p_offset;

			/* update the filesz */
			phdr->p_filesz += sizeof(decoder);
			phdr->p_memsz += sizeof(decoder);

			/* DEBUG usage */
			Elf64_Shdr *shdr = ptr + hdr->e_shoff + 12 * sizeof(Elf64_Shdr);
			shdr->sh_size += sizeof(decoder);

			break;
		}
	}


	int writer = -1;
	if (0 > (writer = open(dst, O_WRONLY | O_CREAT, 0755))) {
		DEBUG("Cannot open file %s", dst);
		return -1;
	}
	write(writer, ptr, total_filesz + sizeof(decoder));

	close(fd);
	close(writer);
	return 0;
}
