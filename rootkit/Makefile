SRC=$(wildcard *.c)
OBJ=$(subst .c,.o,$(SRC))

CFLAGS=-Wall -fPIC
ROOTKIT=rootkit.so

.PHONY: all clean

all: $(ROOTKIT)

$(ROOTKIT): $(OBJ)
	$(CC) -shared $(CFLAGS) $(LDFLAGS) -o $(ROOTKIT) $^

clean: clear

clear:
	rm -rf $(OBJ) $(ROOTKIT)
