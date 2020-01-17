/* Copyright (C) 2020-2020 cmj. All right reserved. */
#include <stdlib.h>

#include "demo.h"
#define PORT	80

int main(int argc, char *argv[]) {
	int nr = 0;

	if (argc < 3) {
		DEBUG("%s IP NR", argv[0]);
		return -1;
	}

	if (0 >= (nr = atoi(argv[2]))) {
		DEBUG("number of threads incorrect : %s", argv[2]);
		return -1;
	}

	DEBUG("Run %s / %d threads", argv[1], nr);
	#pragma omp parallel for
	for (int i = 0; i < nr; ++i) {
		int sk = -1;

		sk = create_client(argv[1], PORT);
		request_handler(sk);
	}

	return 0;
}

