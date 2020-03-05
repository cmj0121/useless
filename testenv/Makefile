.PHONY: all help clean


all: help

help:   # show this help message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

clean:  # remove all test file

%-start: %.yaml
	docker-compose -f $^ up --build -d

%-stop: %.yaml
	docker-compose -f $^ stop

%-restart: %.yaml
	docker-compose -f $^ stop
	docker-compose -f $^ rm -f
	docker-compose -f $^ up --build -d

%-log: %.yaml
	docker-compose -f $^ logs -f
