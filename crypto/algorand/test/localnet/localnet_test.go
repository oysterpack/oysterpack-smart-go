package localnet

import (
	"context"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"testing"
)

func TestGetPrefundedTestAccounts(t *testing.T) {
	accounts := GetPrefundedTestAccounts(t)
	if len(accounts) == 0 {
		t.Fatal("no accounts were returned")
	}
	// accounts should be sorted by ALGO balance desc
	var prevAccount *models.Account
	for _, account := range accounts {
		t.Log(account.Address, account.Amount)
		if prevAccount != nil {
			if account.Amount > prevAccount.Amount {
				t.Error("Accounts are not sorted by ALGO balance DESC", prevAccount.Amount, account.Amount)
			}
		}
		prevAccount = &account
	}

}

func TestIndexerClient(t *testing.T) {
	client := IndexerClient(t)

	// lookup accounts in the indexer
	accounts := GetPrefundedTestAccounts(t)
	for _, account := range accounts {
		validRound, account, err := client.LookupAccountByID(account.Address).Do(context.Background())
		if err != nil {
			t.Error("failed to lookup account")
		}
		t.Log(validRound, account.Address, account.Amount)
	}
}

func TestGenerateTestAccount(t *testing.T) {
	amount := types.ToMicroAlgos(0.1)
	acct := GenerateTestAccount(t, amount)
	algodClient := AlgodClient(t)
	acctInfo, err := algodClient.AccountInformation(acct.Address.String()).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if acctInfo.Amount != uint64(amount) {
		t.Errorf("Account ALGO balance does not match: expected = %v actual = %v", amount, acctInfo.Amount)
	}
}
