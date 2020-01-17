/* Copyright (C) 2020-2020 cmj. All right reserved. */

#include "demo.h"
#define PORT	80

int main(int argc, char *argv[]) {
	if (argc < 2) {
		DEBUG("%s IP", argv[0]);
		return -1;
	}

	int sk = -1;

	sk = create_client(argv[1], PORT);
	request_handler(sk);

	return 0;
}

