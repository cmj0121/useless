.PHONY: all

SRC=$(wildcard web/*)
BUILD=build
RELEASE=release
TIMER=timer

all: $(RELEASE)

help:	# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

$(RELEASE): $(BUILD) $(TIMER)	# build the release code
	mkdir $@
	cp -a build/{index.html,styles.css,main.dart.js} $@/
	cp timer $@

$(BUILD): $(SRC)				# build the web-view
	pub get
	webdev build --release

$(TIMER): $(TIMER).go			# build the server code
	go build -o $@ $^

clean:
	rm -rf $(BUILD) $(RELEASE)
