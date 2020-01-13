.PHONY: all gofmt test

SRC=$(wildcard *.go) $(wildcard */*.go)
DEMO=demo/main.go


all: test gofmt

test: $(DEMO)
	go test -v ./...
	go run $^

gofmt: $(SRC)
	gofmt -w -s $^

.PHONY: clean

clean:
