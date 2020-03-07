.PHONY: help clean gofmt build

SRC=$(wildcard *.go)
BIN=shorturl

help:	# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

build: $(BIN) # build the binary

$(BIN): $(SRC) gofmt
	go build -o $@ $(SRC)

gofmt: $(SRC)
	gofmt -w -s $^

clean:	# cleanup the environment
	rm -f $(BIN)
