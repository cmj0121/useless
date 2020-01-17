
C_SRC=$(wildcard *.c)
BIN=$(subst .c,.example,$(C_SRC))
CFLAGS=-O2 -lpthread -fopenmp -lomp

.PHONY: all

all: $(BIN)

%.example: %.c
	$(CC) $(CFLAGS) -o $@ $^

.PHONY: clean

clean:
	rm -f $(BIN)
