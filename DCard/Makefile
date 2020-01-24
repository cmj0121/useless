SRC=$(wildcard *.go) $(wildcard */*.go)
BIN=dcard

.PHONY: all $(BIN) gofmt

all: $(BIN)

$(BIN): demo/main.go gofmt
	go build -o $@ $<

gofmt: $(SRC)
	gofmt -w -s $^

.PHONY: clean

clean:
	rm -f $(BIN)
