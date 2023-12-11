package account

import (
	"context"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Address is the human-readable Algorand account [address] which maps to the account's public key
//
// [address] = https://developer.algorand.org/docs/get-details/accounts/#transformation-public-key-to-algorand-address
type Address string

// GetAuthAddr looks up the authorized signing account for the specified account address.
//
// If the account is not rekeyed, then the authorized account is itself.
func GetAuthAddr(algodClient *algod.Client, address Address) (authAddr Address, err error) {
	account, err := algodClient.AccountInformation(string(address)).Do(context.Background())
	if err != nil {
		return "", errGetAuthAddrFailed(err)
	}
	if account.AuthAddr == "" {
		return address, nil
	}
	return Address(account.AuthAddr), nil
}

// MakeRekeyTransaction constructs a transaction to rekey the account from the specified address to the specified address
func MakeRekeyTransaction(algodClient *algod.Client, from, to Address) (types.Transaction, error) {
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		return types.Transaction{}, errGetSuggestedParamsFailed(err)
	}
	rekeyTxn, err := transaction.MakePaymentTxn(string(from), string(from), 0, nil, "", sp)
	if err != nil {
		return types.Transaction{}, errMakePaymentTxn(err)
	}
	if err = rekeyTxn.Rekey(string(to)); err != nil {
		return types.Transaction{}, errSettingRekeyTo(to, err)
	}
	return rekeyTxn, nil
}
