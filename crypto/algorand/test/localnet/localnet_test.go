package localnet

import (
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"testing"
)

func TestGetPrefundedTestAccounts(t *testing.T) {
	accounts := GetPrefundedTestAccounts(t)
	if len(accounts) == 0 {
		t.Error("no accounts were returned")
		t.FailNow()
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
