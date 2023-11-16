package kmd

import (
	"errors"
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"strings"
)

const walletDriverName = "sqlite"

type Wallet struct {
	Id   string
	Name string
}

// WalletManager
//
// [KMD]: https://developer.algorand.org/docs/get-details/accounts/create/#wallet-derived-kmd
type WalletManager interface {
	// List returns the list of wallets that are being managed
	List() ([]Wallet, error)

	// Create a new wallet using the specified name and password
	Create(name, password string) (Wallet, error)

	Recover(name, password, backupPhrase string) (Wallet, error)

	// ExportbackupPhrase exports the mnemonic for the wallet's master derivation key
	ExportBackupPhrase(name, password string) (string, error)
}

// KMD based implementation for WalletManager
type kmdWalletManager struct {
	kmdClient kmd.Client
}

// New constructs a new WalletManager instance using the specified KMD client
//
// The WalletManager is backed by a single KMD instance
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
// - name and password whitespace will be trimmed
//
// ## Constraints
// - name and password must not be blank
// - wallet with the same name must not already exist
func (walletManager *kmdWalletManager) Create(name, password string) (Wallet, error) {
	name, password, err := trimNamePassword(name, password)
	if err != nil {
		return Wallet{}, err
	}

	response, err := walletManager.kmdClient.CreateWallet(
		name,
		password,
		walletDriverName,
		types.MasterDerivationKey{},
	)
	if err != nil {
		return Wallet{}, err
	}
	wallet := Wallet{
		Id:   response.Wallet.ID,
		Name: response.Wallet.Name,
	}
	return wallet, nil
}

func trimNamePassword(name, password string) (walletName string, walletPassword string, err error) {
	walletName = strings.TrimSpace(name)
	if len(walletName) == 0 {
		return walletName, password, errors.New("name cannot be blank")
	}
	walletPassword = strings.TrimSpace(password)
	if len(walletPassword) == 0 {
		return walletName, walletPassword, errors.New("password cannot be blank")
	}
	return
}

func (walletManager *kmdWalletManager) wallet(name string) (Wallet, error) {
	wallets, err := walletManager.List()
	if err != nil {
		return Wallet{}, err
	}
	for _, wallet := range wallets {
		if wallet.Name == name {
			return wallet, nil
		}
	}

	return Wallet{}, errors.New("wallet does not exist")
}

func (walletManager *kmdWalletManager) walletHandle(name, password string) (handle string, err error) {
	wallet, err := walletManager.wallet(name)
	if err != nil {
		return "", err
	}

	response, err := walletManager.kmdClient.InitWalletHandle(wallet.Id, password)
	if err != nil {
		return "", err
	}
	return response.WalletHandleToken, nil
}

func (walletManager *kmdWalletManager) ExportBackupPhrase(name, password string) (backupPhrase string, err error) {
	name, password, err = trimNamePassword(name, password)
	if err != nil {
		return "", err
	}

	handle, err := walletManager.walletHandle(name, password)
	if err != nil {
		return "", err
	}

	response, err := walletManager.kmdClient.ExportMasterDerivationKey(handle, password)
	if err != nil {
		return "", err
	}

	return mnemonic.FromMasterDerivationKey(response.MasterDerivationKey)
}

func (walletManager *kmdWalletManager) Recover(name, password, backupPhrase string) (Wallet, error) {
	name, password, err := trimNamePassword(name, password)
	if err != nil {
		return Wallet{}, err
	}

	keyBytes, err := mnemonic.ToKey(backupPhrase)
	var mdk types.MasterDerivationKey
	copy(mdk[:], keyBytes)
	if err != nil {
		return Wallet{}, err
	}

	response, err := walletManager.kmdClient.CreateWallet(
		name,
		password,
		walletDriverName,
		mdk,
	)
	if err != nil {
		return Wallet{}, err
	}
	wallet := Wallet{
		Id:   response.Wallet.ID,
		Name: response.Wallet.Name,
	}
	return wallet, nil
}
