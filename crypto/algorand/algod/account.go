package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
)

// GetAuthAddr looks up the authorized signing account for the specified account address.
//
// If the account is not rekeyed, then the authorized account is itself.
func GetAuthAddr(algodClient *algod.Client, address string) (authAddr string, err error) {
	account, err := algodClient.AccountInformation(address).Do(context.Background())
	if err != nil {
		return "", err
	}
	if account.AuthAddr == "" {
		return address, nil
	}
	return account.AuthAddr, nil
}
