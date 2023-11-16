pre-commit: format localnet-reset test

test:
	cd crypto/algorand/kmd && go test


format:
	cd crypto/algorand/kmd && go fmt

godoc:

# algorand node commands depend on the $ALGORAND_DATA env var
# if not set, then it defaults to /var/lib/algorand
check_algorand_data_env_var:
ALGORAND_DATA ?= "/var/lib/algorand"

algod-status: check_algorand_data_env_var
	sudo -u algorand goal -d $(ALGORAND_DATA) node status

kmd-start: check_algorand_data_env_var
	sudo -u algorand goal -d $(ALGORAND_DATA) kmd start

kmd-stop: check_algorand_data_env_var
	sudo -u algorand goal -d $(ALGORAND_DATA) kmd stop

localnet-reset:
	algokit localnet reset