/* Copyright (C) 2020-2020 cmj. All right reserved. */
#include <stdlib.h>
#include <pthread.h>

#include "demo.h"
#define PORT	80

void* job(void *ip) {
	int sk = -1;

	sk = create_client(ip, PORT);
	request_handler(sk);
	return NULL;
}

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

	pthread_t *pids = NULL;
	if (NULL == (pids = malloc(sizeof(pthread_t) * nr))) {
		DEBUG("Cannot malloc %d pthread_t", nr);
		return -1;
	}

	for (int i = 0; i < nr; ++i) {
		pthread_create(&pids[i], NULL, job, argv[1]);
	}

	for (int i = 0; i < nr; ++i) {
		pthread_join(pids[i], NULL);
	}


	free(pids);
	return 0;
}

