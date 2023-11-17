package kmd

import (
	"errors"
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
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

	// Contains reports true if the wallet for the specified name exists
	Contains(name string) (bool, error)

	// Create a new wallet using the specified name and password
	Create(name, password string) (Wallet, error)

	Recover(name, password, backupPhrase string) (Wallet, error)

	// ExportBackupPhrase exports the mnemonic for the wallet's master derivation key
	ExportBackupPhrase(name, password string) (string, error)

	// ListAccounts returns the account addresses for the specified wallet
	ListAccounts(name, password string) (addresses []string, err error)

	// ContainsAccount is used to check if an address exists within the specified wallet
	ContainsAccount(name, password, address string) (bool, error)

	// CreateAccount generates a new account within the specified wallet
	CreateAccount(name, password string) (address string, err error)

	// CreateAccounts will attempt to create the specified number of accounts.
	//
	// If an error occurs while generating accounts, then the error will be returned along with the accounts that
	// were created up to that point. Thus, even if an error is returned, check if any addresses were created.
	CreateAccounts(name, password string, count uint) (addresses []string, err error)

	// DeleteAccount deletes the specified address located within the specified wallet.
	//
	// It's ok to try to delete an account that does exist, i.e., no error is returned
	DeleteAccount(name, password, address string) error

	// DeleteAccounts deletes the specified list accounts for the specified wallet.
	//
	// - At least 1 account address to delete must be specified.
	//
	// Any errors that occur while deleting an account are returned in the map.
	// If there is no entry for the address in the returned map, then it means the account was successfully deleted.
	// If all deletes completed successfully, then nil will be returned.
	DeleteAccounts(name, password string, addresses ...string) (deleteErrors map[string]error, err error)

	// ExportPrivateKey exports the account's private key in [mnemonic] form
	//
	// [mnemonic]: https://developer.algorand.org/docs/get-details/accounts/#transformation-private-key-to-25-word-mnemonic
	ExportPrivateKey(name, password, address string) (passPhrase string, err error)
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

func (walletManager *kmdWalletManager) Contains(name string) (bool, error) {
	kmdWallets, err := walletManager.kmdClient.ListWallets()
	if err != nil {
		return false, err
	}
	for _, wallet := range kmdWallets.Wallets {
		if wallet.Name == name {
			return true, nil
		}
	}

	return false, nil
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

	walletHandle, err := walletManager.walletHandle(name, password)
	if err != nil {
		return "", err
	}
	defer walletManager.releaseWalletHandle(name, walletHandle)

	response, err := walletManager.kmdClient.ExportMasterDerivationKey(walletHandle, password)
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
	defer walletManager.releaseWalletHandle(name, walletHandle)

	listKeysResponse, err := walletManager.kmdClient.ListKeys(walletHandle)
	if err != nil {
		return nil, err
	}
	return listKeysResponse.Addresses, nil
}

func (walletManager *kmdWalletManager) ContainsAccount(name, password, address string) (bool, error) {
	walletAddresses, err := walletManager.ListAccounts(name, password)
	if err != nil {
		return false, err
	}
	for _, walletAddress := range walletAddresses {
		if walletAddress == address {
			return true, nil
		}
	}
	return false, nil
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
	defer walletManager.releaseWalletHandle(name, walletHandle)

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
	defer walletManager.releaseWalletHandle(name, walletHandle)

	for i := 0; uint(i) < count; i++ {
		response, err := walletManager.kmdClient.GenerateKey(walletHandle)
		if err != nil {
			return addresses, err
		}
		addresses = append(addresses, response.Address)
	}

	return addresses, nil
}

func (walletManager *kmdWalletManager) DeleteAccount(name, password, address string) error {
	name, password, err := trimNamePassword(name, password)
	if err != nil {
		return err
	}

	walletHandle, err := walletManager.walletHandle(name, password)
	if err != nil {
		return err
	}
	defer walletManager.releaseWalletHandle(name, walletHandle)

	_, err = walletManager.kmdClient.DeleteKey(walletHandle, password, address)
	return err
}

func (walletManager *kmdWalletManager) DeleteAccounts(name, password string, addresses ...string) (deleteErrors map[string]error, err error) {
	if len(addresses) == 0 {
		return nil, errors.New("at least 1 account to delete must be specified")
	}

	name, password, err = trimNamePassword(name, password)
	if err != nil {
		return nil, err
	}

	walletHandle, err := walletManager.walletHandle(name, password)
	if err != nil {
		return nil, err
	}
	defer walletManager.releaseWalletHandle(name, walletHandle)

	deleteErrors = make(map[string]error)
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(addresses))

	for _, address := range addresses {
		go func(address string) {
			defer wg.Done()
			_, err = walletManager.kmdClient.DeleteKey(walletHandle, password, address)
			if err != nil {
				lock.Lock()
				deleteErrors[address] = err
				lock.Unlock()
			}
		}(address)
	}

	wg.Wait()

	if len(deleteErrors) == 0 {
		return nil, nil
	}
	return
}

func (walletManager *kmdWalletManager) ExportPrivateKey(name, password, address string) (passPhrase string, err error) {
	name, password, err = trimNamePassword(name, password)
	if err != nil {
		return "", err
	}

	walletHandle, err := walletManager.walletHandle(name, password)
	if err != nil {
		return "", err
	}
	defer walletManager.releaseWalletHandle(name, walletHandle)

	response, err := walletManager.kmdClient.ExportKey(walletHandle, password, address)
	if err != nil {
		return "", err
	}
	return mnemonic.FromPrivateKey(response.PrivateKey)
}

func (walletManager *kmdWalletManager) releaseWalletHandle(walletName, walletHandle string) {
	_, err := walletManager.kmdClient.ReleaseWalletHandle(walletHandle)
	if err != nil {
		log.Warning("Failed to release KMD handle for wallet:", walletName)
	}
}
