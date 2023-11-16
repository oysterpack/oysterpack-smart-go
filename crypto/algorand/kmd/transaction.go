package kmd

import (
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/oysterpack/oysterpack-smart-go/crypto/algorand/model"
)

// KMD wallet ID
type WalletId string

// KMD wallet account
type WalletAccount struct {
	WalletId
	model.Address
}

// WalletSession is a connected KMD wallet. The wallet will automatically handle
// renewing its wallet handle in the background.
//
// WalletSession is also a TransactionSigner and knows how to handle rekeyed accounts. If the transaction sender
// account has been rekeyed, then the authorized account will be used to sign the transaction. The authorized account
// must exist in this wallet.
type Signer interface {
	transaction.TransactionSigner
}

type KmdService struct {
	kmd.Client
}
