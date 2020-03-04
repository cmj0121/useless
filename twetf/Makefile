.PHONY: all help clean gofmt
SRC=$(wildcard *.go)
BIN=twetf

all: $(BIN)	# default action

help:		# show this help message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

clean:		# cleanup environment
	rm -f $(BIN)

gofmt: $(SRC)
	@gofmt -s -w $(SRC)

$(BIN): $(SRC)
	@go build -o $@ $^
