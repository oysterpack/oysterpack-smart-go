// Package localnet provides access to the localnet [AlgoKit] environment.
//
// [AlgoKit]: https://developer.algorand.org/docs/get-details/algokit/

package localnet

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
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

func GetDefaultWalletHandle(kmdClient kmd.Client) (string, error) {
	listWalletsResponse, err := kmdClient.ListWallets()
	if err != nil {
		return "", err
	}
	for _, wallet := range listWalletsResponse.Wallets {
		if wallet.Name == DefaultKmdWalletName {
			initWalletHandleResponse, err := kmdClient.InitWalletHandle(wallet.ID, DefaultKmdWalletPassword)
			if err != nil {
				return "", fmt.Errorf("failed to init wallet handle for default wallet : %w", err)
			}
			walletHandle := initWalletHandleResponse.WalletHandleToken
			return walletHandle, nil
		}
	}

	return "", fmt.Errorf("default wallet not not found: %v", DefaultKmdWalletName)
}

func ReleaseWalletHandle(kmdClient kmd.Client, walletHandle string) {
	_, err := kmdClient.ReleaseWalletHandle(walletHandle)
	if err != nil {
		panic(fmt.Errorf("failed to release wallet handle: %w", err))
	}
}

// GetPrefundedTestAccounts returns accounts from the default KMD wallet sorted by ALGO account balance in descending order,
// i.e., the account with the highest balance is first.
func GetPrefundedTestAccounts(t *testing.T) []models.Account {
	kmdClient := KmdClient(t)
	algodClient := AlgodClient(t)

	defaultWalletHandle, err := GetDefaultWalletHandle(kmdClient)
	if err != nil {
		t.Fatal("failed to load prefunded accounts from KMD default wallet:", err)
	}
	defer ReleaseWalletHandle(kmdClient, defaultWalletHandle)
	listKeysResponse, err := kmdClient.ListKeys(defaultWalletHandle)
	if err != nil {
		t.Fatal("failed to list accounts for default wallet", err)
	}
	addresses := listKeysResponse.Addresses
	accounts := make([]models.Account, 0)
	for _, address := range addresses {
		account, err := algodClient.AccountInformation(address).Do(context.Background())
		if err != nil {
			t.Fatal("failed to retrieve account info", err)
		}
		if account.Amount > 0 {
			accounts = append(accounts, account)
		}
	}

	if len(accounts) == 0 {
		t.Fatal("No prefunded accounts were found")
	}

	// sort the accounts by ALGO account balance in descending order
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].Amount > accounts[j].Amount
	})
	return accounts
}

// GenerateTestAccount generates a new test Algorand account and funds the account with the specified ALGO amount
//   - The account is funded from one of the prefunded accounts that has the most ALGO. The ALGO is transferred from
//     the prefunded account to the new generated account
func GenerateTestAccount(t *testing.T, algos types.MicroAlgos) crypto.Account {
	account := crypto.GenerateAccount()

	algodClient := AlgodClient(t)
	if algos == 0 {
		t.Fatal("account must be funded")
	}

	fundingAccount := GetPrefundedTestAccounts(t)[0]

	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		t.Fatalf("error getting suggested tx params: %s", err)
	}

	txn, err := transaction.MakePaymentTxn(
		fundingAccount.Address,
		account.Address.String(), uint64(algos),
		nil,
		"",
		sp,
	)
	if err != nil {
		t.Fatalf("failed to make transaction: %s", err)
	}

	kmdClient := KmdClient(t)
	defaultWalletHandle, err := GetDefaultWalletHandle(kmdClient)
	if err != nil {
		t.Fatal(err)
	}
	defer ReleaseWalletHandle(kmdClient, defaultWalletHandle)

	resp, err := kmdClient.SignTransaction(defaultWalletHandle, DefaultKmdWalletPassword, txn)
	if err != nil {
		t.Fatal(err)
	}
	txID, err := algodClient.SendRawTransaction(resp.SignedTransaction).Do(context.Background())
	if err != nil {
		t.Fatal("failed to send transaction:", err)
	}
	_, err = transaction.WaitForConfirmation(algodClient, txID, 4, context.Background())
	if err != nil {
		t.Fatalf("Error waiting for confirmation on txID: %s\n", txID)
	}

	return account
}
