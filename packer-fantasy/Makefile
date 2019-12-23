SRC=$(wildcard *.go)

SUBDIR=

.PHONY: all gofmt test clean demo

all: gofmt test demo

gofmt:
	gofmt -s -w $(SRC)

test: gofmt
	go test -v ./...

demo: test
	go build demo/main.go

.PHONY: clean

clean:
