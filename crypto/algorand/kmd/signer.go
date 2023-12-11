package kmd

import (
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
)

// TransactionSigner implements the algorand [transaction.TransactionSigner] interface
//
// [transaction.TransactionSigner] = https://pkg.go.dev/github.com/algorand/go-algorand-sdk/v2@v2.3.0/transaction#TransactionSigner
type TransactionSigner struct {
	kmdClient      *kmd.Client
	walletName     string
	walletPassword string
}

//func (signer TransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
//
//	for _, indexToSign := range indexesToSign {
//		txn := txGroup[indexToSign]
//	}
//}
