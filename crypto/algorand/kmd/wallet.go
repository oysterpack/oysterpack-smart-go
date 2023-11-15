package kmd

import "github.com/algorand/go-algorand-sdk/v2/client/kmd"

type Wallet struct {
	Id   string
	Name string
}

// WalletManager
//
// [KMD]: https://developer.algorand.org/docs/get-details/accounts/create/#wallet-derived-kmd
type WalletManager interface {
	List() ([]Wallet, error)
}

// KMD based implementation for WalletManager
type kmdWalletManager struct {
	kmdClient kmd.Client
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

// New constructs a new WalletManager instance using the specified KMD client
func New(kmdClient kmd.Client) WalletManager {
	return &kmdWalletManager{
		kmdClient: kmdClient,
	}
}
