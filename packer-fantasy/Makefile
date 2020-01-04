SRC=$(wildcard *.c)
BIN=pf

SUBDIR=

.PHONY: all

all: $(BIN)

$(BIN): $(SRC)
	gcc -o $@ $^

.PHONY: clean

clean:
	rm -f $(BIN)
