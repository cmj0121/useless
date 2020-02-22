GO-SRC=$(wildcard *.go)

.PHONY: all run gofmt

all: gofmt

run: gofmt
	go run lazydash.go

gofmt: $(GO-SRC)
	gofmt -s -w $^

.PHONY: clean

clean:
	find . -name '*.swp' -delete
