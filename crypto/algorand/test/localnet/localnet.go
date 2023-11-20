// Package localnet provides access to the localnet [AlgoKit] environment.
//
// [AlgoKit]: https://developer.algorand.org/docs/get-details/algokit/

package localnet

import (
	"context"
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
	"sort"
	"strings"
	"testing"
)

const DefaultKmdWalletName = "unencrypted-default-wallet"
const DefaultKmdWalletPassword = ""

func KmdClient(t *testing.T) kmd.Client {
	const url = "http://localhost:4002"
	var token = strings.Repeat("a", 64)

	client, err := kmd.MakeClient(
		url,
		token,
	)

	if err != nil {
		t.Fatal("failed to connect to KMD on localnet", err)
	}
	return client
}

func AlgodClient(t *testing.T) *algod.Client {
	const url = "http://localhost:4001"
	var token = strings.Repeat("a", 64)

	client, err := algod.MakeClient(url, token)
	if err != nil {
		t.Fatal("failed to connect to algod on localnet", err)
	}

	return client
}

func IndexerClient(t *testing.T) *indexer.Client {
	const url = "http://localhost:8980"
	var token = strings.Repeat("a", 64)

	client, err := indexer.MakeClient(url, token)
	if err != nil {
		t.Fatal("failed to connect to indexer on localnet", err)
	}

	return client
}

// GetPrefundedTestAccounts returns accounts from the default KMD wallet sorted by ALGO account balance in descending order,
// i.e., the account with the highest balance is first.
func GetPrefundedTestAccounts(t *testing.T) []models.Account {
	kmdClient := KmdClient(t)
	algodClient := AlgodClient(t)

	listWalletsResponse, err := kmdClient.ListWallets()
	if err != nil {
		t.Fatal("failed to list wallets", err)
	}
	for _, wallet := range listWalletsResponse.Wallets {
		if wallet.Name == DefaultKmdWalletName {
			initWalletHandleResponse, err := kmdClient.InitWalletHandle(wallet.ID, DefaultKmdWalletPassword)
			if err != nil {
				t.Fatal("failed to init wallet handle for default wallet")
			}
			walletHandle := initWalletHandleResponse.WalletHandleToken
			defer func() {
				_, err := kmdClient.ReleaseWalletHandle(walletHandle)
				if err != nil {
					t.Error("failed to release wallet handle", err)
				}
			}()
			listKeysResponse, err := kmdClient.ListKeys(walletHandle)
			if err != nil {
				t.Fatal("failed to list accounts for default wallet", err)
			}
			addresses := listKeysResponse.Addresses
			accounts := make([]models.Account, 0)
			for _, address := range addresses {
				account, err := AccountInfo(algodClient, address)
				if err != nil {
					t.Fatal("failed to retrieve account info", err)
				}
				if account.Amount > 0 {
					accounts = append(accounts, account)
				}
			}

			if len(accounts) == 0 {
				t.Fatal("No prefunded accounts were found")
			} else {
				// sort the accounts by ALGO account balance in descending order
				sort.Slice(accounts, func(i, j int) bool {
					return accounts[i].Amount > accounts[j].Amount
				})
				return accounts
			}
		}
	}
	panic("failed to load prefunded accounts from KMD default wallet")
}

func AccountInfo(algodClient *algod.Client, address string) (models.Account, error) {
	account, err := algodClient.AccountInformation(address).Do(context.Background())
	if err != nil {
		return models.Account{}, err
	}
	return account, nil
}
