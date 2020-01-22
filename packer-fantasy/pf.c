#include <stdio.h>
#include <stdlib.h>
#include <elf.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>

#define DEBUG(msg, ...)	\
	fprintf(stderr, "[%s #%04d] " msg "\n", __FILE__, __LINE__, ##__VA_ARGS__); \

off_t Shift(void *bin, off_t len, unsigned char key) {
	Elf64_Ehdr *hdr = bin;
	off_t entry = hdr->e_entry;
	off_t code_off = 0;
	off_t code_len = 0;

	/* Find entry pointer */
	DEBUG("Entry Point : vaddr 0x%lX", entry);

	/* Find machine code */
	off_t offset = hdr->e_phoff;
	for (int i = 0; i < hdr->e_phnum; ++i) {
		Elf64_Phdr *phdr = bin + offset;

		if (phdr->p_type == PT_LOAD) {
			phdr->p_flags = PF_R | PF_X | PF_W;
			DEBUG("Find LOAD program on #%d, change to PF_R|PF_W|PF_X, change to PF_R|PF_W|PF_X", i);

			code_off = entry - phdr->p_vaddr;
			code_len = phdr->p_filesz;
			DEBUG("Code located on file offset 0x%lx, %ld bytes", code_off, code_len);
			break;
		}
	}

	/* Encode machine code */
	unsigned char enc_key = key;
	unsigned char *code = bin + code_off;
	for (off_t i = 0; i < code_len; ++i) {
		code[i] = code[i] ^ enc_key;
		enc_key ++;
	}

	/* Add decoder */
	unsigned char decoder[] = {
		/* 0000 mov rax EntryPoint */
		0xB8, 0xAA, 0xAA, 0xAA, 0xAA,
		/* 0005 mov rcx CodeLength */
		0xB9, 0xBB, 0xBB, 0xBB, 0xBB,
		/* 000A mov dl KEY */
		0xB2, 0xCC,
		/* 000C xor byte [rax] dl */
		0x30, 0x10,
		/* 000E inc dl */
		0xFE, 0xC2,
		/* 0010 inc rax */
		0x48, 0xFF, 0xC0,
		/* 0013 loop */
		0xE2, 0xF7,
		/* 0015 jmp EntryPoint */
	 	0xE9, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
	};
	uint8_t *ptr8 = NULL;
	uint32_t *ptr32 = NULL;

	/* setup the entry pointer */
	ptr32 = (uint32_t *)(decoder + 0x1);
	*ptr32 = (uint32_t)entry;
	/* setup the code length  */
	ptr32 = (uint32_t *)(decoder + 0x6);
	*ptr32 = (uint32_t)code_len;
	/* set the key */
	ptr8 = (uint8_t *)(decoder + 0x0B);
	*ptr8 = (uint8_t)key;
	/* setup the original address  */
	ptr32 = (uint32_t *)(decoder + 0x16);
	*ptr32 = (uint32_t)(0) - (uint32_t)(code_len + sizeof(decoder)) + 1;

	for (int i = 0; i < sizeof(decoder); ++i) {
		code[code_len + i] = decoder[i];
	}

	/* Update entry pointer */
	hdr->e_entry = entry + code_len;
	DEBUG("Change to new entry addr : vaddr 0x%lX", hdr->e_entry);
	return len + sizeof(decoder);
}


int main(int argc, char* argv[]) {
	int fd = -1;

	if (argc < 2) {
		DEBUG("%s FILE", argv[0]);
		return -1;
	}

	DEBUG("WARNING: !! This is the experience tool and only workable on ELF / 1 LOAD binary !!");

	if (0 > (fd = open(argv[1], O_RDONLY))) {
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
	if (hdr->e_ident[EI_MAG0] != ELFMAG0 || hdr->e_ident[EI_MAG1] != ELFMAG1 || hdr->e_ident[EI_MAG2] != ELFMAG2) {
		DEBUG("Not the ELF file");
		return -1;
	} else if (hdr->e_ident[EI_CLASS] != ELFCLASS64) {
		DEBUG("Not the ELF/64 file");
		return -1;
	}

	/* packer */
	total_filesz = Shift(ptr, total_filesz, 0x12);

	int writer = -1;
	char dst[] = "packed";
	if (0 > (writer = open(dst, O_WRONLY | O_CREAT, 0755))) {
		DEBUG("Cannot open file %s", dst);
		return -1;
	}
	write(writer, ptr, total_filesz);

	close(fd);
	close(writer);
	return 0;
}
