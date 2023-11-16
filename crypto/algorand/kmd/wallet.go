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

	// Get Wallet by name
	//
	// If the wallet is not found, then an error is returned
	Get(name string) (Wallet, error)

	// Create a new wallet using the specified name and password
	Create(name, password string) (Wallet, error)

	Recover(name, password, backupPhrase string) (Wallet, error)

	// ExportBackupPhrase exports the mnemonic for the wallet's master derivation key
	ExportBackupPhrase(name, password string) (string, error)

	// ListAccounts returns the account addresses for the specified wallet
	ListAccounts(name, password string) (addresses []string, err error)

	CreateAccount(name, password string) (address string, err error)

	// CreateAccounts will attempt to create the specified number of accounts.
	//
	// If an error occurs while generating accounts, then the error will be returned along with the accounts that
	// were created up to that point. Thus, even if an error is returned, check if any addresses were created.
	CreateAccounts(name, password string, count uint) (addresses []string, err error)
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
		return nil, err
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

func (walletManager *kmdWalletManager) Get(name string) (Wallet, error) {
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
	wallet, err := walletManager.Get(name)
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
	defer func() {
		_, err = walletManager.kmdClient.ReleaseWalletHandle(handle)
	}()

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

func (walletManager *kmdWalletManager) ListAccounts(name, password string) (addresses []string, err error) {
	name, password, err = trimNamePassword(name, password)
	if err != nil {
		return nil, err
	}

	wallet, err := walletManager.Get(name)
	if err != nil {
		return nil, err
	}

	walletHandleResponse, err := walletManager.kmdClient.InitWalletHandle(wallet.Id, password)
	if err != nil {
		return nil, err
	}
	walletHandle := walletHandleResponse.WalletHandleToken
	defer func() {
		_, err = walletManager.kmdClient.ReleaseWalletHandle(walletHandle)
	}()

	listKeysResponse, err := walletManager.kmdClient.ListKeys(walletHandle)
	if err != nil {
		return nil, err
	}
	return listKeysResponse.Addresses, nil
}

func (walletManager *kmdWalletManager) CreateAccount(name, password string) (address string, err error) {
	name, password, err = trimNamePassword(name, password)
	if err != nil {
		return "", err
	}

	walletHandle, err := walletManager.walletHandle(name, password)
	if err != nil {
		return "", err
	}
	defer func() {
		_, err = walletManager.kmdClient.ReleaseWalletHandle(walletHandle)
	}()

	response, err := walletManager.kmdClient.GenerateKey(walletHandle)
	if err != nil {
		return "", err
	}
	return response.Address, nil
}

func (walletManager *kmdWalletManager) CreateAccounts(name, password string, count uint) (addresses []string, err error) {
	name, password, err = trimNamePassword(name, password)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("`count` must be greater than 0")
	}

	walletHandle, err := walletManager.walletHandle(name, password)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, err = walletManager.kmdClient.ReleaseWalletHandle(walletHandle)
	}()

	for i := 0; uint(i) < count; i++ {
		response, err := walletManager.kmdClient.GenerateKey(walletHandle)
		if err != nil {
			return addresses, err
		}
		addresses = append(addresses, response.Address)
	}

	return addresses, nil
}
