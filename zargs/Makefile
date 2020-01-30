SRC=$(wildcard *.go)
DEMO_SRC=demo/main.go

.PHONY: all gofmt demo

all: test demo

gofmt: $(SRC) $(DEMO_SRC)
	gofmt -s -w $^

demo: $(DEMO_SRC)
	go run $^

test: gofmt
	go test -v ./...

.PHONY: clean

clean:
