package kmd

import (
	"github.com/oysterpack/oysterpack-smart-go/algorand/model"
)

type Rekeying interface {
	// Rekey the account
	//
	// Both accounts must exist in
	Rekey(from model.Address, to model.Address) error

	// If the account is rekeyed, then it will rekey the account back to itself
	RekeyBack(account model.Address) error
}
