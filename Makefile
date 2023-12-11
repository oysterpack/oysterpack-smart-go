pre-commit: format lint localnet-reset test

test:
	cd fxapp && go test

format:
	cd fxapp && go fmt

godoc:
	cd crypto/algorand && godoc -http :6060 -index

# runs go vet and staticcheck lint tools on all module packages
lint:
	cd fxapp && go vet ./...
	cd fxapp && staticcheck ./...

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

# resets the local Algorand environment provided by AlgoKit
localnet-reset:
	algokit localnet reset --update