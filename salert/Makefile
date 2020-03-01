SRC=salert
BIN=/usr/local/sbin/salert

.PHONY: all install uninstall

all:	# show this help message
	@printf "Usage: make OPTION \n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

install: $(BIN)	# install the script to local machine

$(BIN): $(SRC)
	install -m750 $(SRC) $(BIN)
	@echo sed -i '/$(SRC)/d'/etc/pam.d/sshd
	@echo "account   optional  pam_exec.so          $(BIN)" >> /etc/pam.d/sshd

uninstall:	# uninstall existed script
	rm -f $(BIN)
	@sed -i "/$(SRC)/d" /etc/pam.d/sshd
