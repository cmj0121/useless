/* Copyright (C) 2020-2020 cmj. All right reserved. */
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include <unistd.h>
#include <sys/syscall.h>
#define gettid() syscall(SYS_gettid)

#define DEBUG(msg, ...)	\
	fprintf(stderr, "[%s L#%04d] (%d/%lld) " msg "\n", __FILE__, __LINE__, getpid(), gettid(), ##__VA_ARGS__);

void request_handler(int sk) {
	if (sk < 0) {
		DEBUG("Invalid socket %d", sk);
		return;
	}

	/* operations : send HTTP GET request */
	const char msg[] = "GET / HTTP/1.1\r\n\r\n";
	ssize_t send_len = 0;
	while (send_len < sizeof(msg)) {
		ssize_t len = 0;

		if (0 > (len =send(sk, (void *)msg + send_len, sizeof(msg) - send_len, 0))) {
			DEBUG("Send fail (%s)", strerror(errno));
			goto END;
		}


		DEBUG("Send %ld bytes", len);
		send_len += len;
	}

	char buff[BUFSIZ] = {0};
	ssize_t recv_len = 0;

	recv_len = recv(sk, buff, sizeof(buff), 0);
	DEBUG("Receive %ld bytes", recv_len);
END:
	close(sk);
	return;
}

int create_client(const char *ip, int port) {
	int sk = -1;
	struct sockaddr_in addr = {0};

	/* Setup the client information */
	addr.sin_family = AF_INET;
	addr.sin_addr.s_addr = inet_addr(ip);
	addr.sin_port = htons(port);	/* convert to network endian */

	DEBUG("Try to connect %s:%d ...", ip, port);
	/* Create socket */
	if (0 > (sk = socket(AF_INET, SOCK_STREAM, 0))) {
		DEBUG("Cannot create socket (%s)", strerror(errno));
		return sk;
	} else if (0 > connect(sk, (struct sockaddr *)&addr, sizeof(addr))) {
		DEBUG("Cannot connect to %s:%d (%s)", ip, port, strerror(errno));
		close(sk);
		return -1;
	}

	return sk;
}
