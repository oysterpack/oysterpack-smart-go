package kmd

import (
	"errors"
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"strings"
)

type Wallet struct {
	Id   string
	Name string
}

// WalletManager
//
// [KMD]: https://developer.algorand.org/docs/get-details/accounts/create/#wallet-derived-kmd
type WalletManager interface {
	List() ([]Wallet, error)

	Create(name, password string) (Wallet, error)
}

// KMD based implementation for WalletManager
type kmdWalletManager struct {
	kmdClient kmd.Client
}

// New constructs a new WalletManager instance using the specified KMD client
func New(kmdClient kmd.Client) WalletManager {
	return &kmdWalletManager{
		kmdClient: kmdClient,
	}
}

// List implements the WalletManager interface
func (walletManager *kmdWalletManager) List() ([]Wallet, error) {
	kmdWallets, err := walletManager.kmdClient.ListWallets()
	if err != nil {
		return []Wallet{}, err
	}
	wallets := make([]Wallet, len(kmdWallets.Wallets))
	for i, wallet := range kmdWallets.Wallets {
		wallets[i] = Wallet{wallet.ID, wallet.Name}
	}

	return wallets, nil
}

// Create a new Wallet with the specified name and password
//
// ## Constraints
// - name and password must not be blank
// - wallet with the same name must not already exist
func (walletManager *kmdWalletManager) Create(name, password string) (Wallet, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return Wallet{}, errors.New("name cannot be blank")
	}
	password = strings.TrimSpace(password)
	if len(password) == 0 {
		return Wallet{}, errors.New("password cannot be blank")
	}

	response, err := walletManager.kmdClient.CreateWallet(
		name,
		password,
		"sqlite",
		types.MasterDerivationKey{},
	)
	if err != nil {
		return Wallet{}, err
	}
	wallet := Wallet{response.Wallet.ID, response.Wallet.Name}
	return wallet, nil
}
