SRC=$(wildcard *.c)
BIN=pf

SUBDIR=

.PHONY: all demo

all: $(BIN)

$(BIN): $(SRC)
	gcc -o $@ $^

demo: $(BIN)
	./$(BIN) demo


.PHONY: clean

clean:
	rm -f $(BIN)
