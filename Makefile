pre-commit: format localnet-reset test

test:
	cd crypto/algorand/kmd && go test
	cd crypto/algorand/test/localnet && go test

format:
	cd crypto/algorand/kmd && go fmt
	cd crypto/algorand/test/localnet && go fmt
	cd crypto/algorand/transaction && go fmt

godoc:
	cd crypto/algorand && godoc -http :6060 -index

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