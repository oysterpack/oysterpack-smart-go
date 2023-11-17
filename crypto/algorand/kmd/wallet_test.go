package kmd_test

import (
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/oklog/ulid/v2"
	"github.com/oysterpack/oysterpack-smart-go/crypto/algorand/kmd"
	"github.com/oysterpack/oysterpack-smart-go/crypto/algorand/test/localnet"
	"sort"
	"strings"
	"testing"
)

func TestListWallets(t *testing.T) {
	kmdClient := localnet.KmdClient(t)

	// Get list of wallets from the KMD service directly
	kmdWallets, err := kmdClient.ListWallets()
	if err != nil {
		t.Log("Failed to list walletManager", err)
		t.FailNow()
	}
	t.Logf("KMD Wallet count = %v", len(kmdWallets.Wallets))
	for _, wallet := range kmdWallets.Wallets {
		t.Log(wallet)
	}

	// Create a WalletManager instance
	walletManager := kmd.New(kmdClient)
	// Get the list of wallets using the WalletManager
	wallets, err := walletManager.List()
	if err != nil {
		t.Error("Failed to list walletManager", err)
		return
	}
	t.Log(wallets)
	// Verify that the number of wallets matches
	if len(wallets) != len(kmdWallets.Wallets) {
		t.Errorf("Wallet count did not match: expected=%v actual=%v", len(kmdWallets.Wallets), len(wallets))
	}

	// Verify that the walletManager match

	// First sort the walletManager by ID
	sort.Slice(kmdWallets.Wallets, func(i, j int) bool {
		return kmdWallets.Wallets[i].ID < kmdWallets.Wallets[j].ID
	})
	sort.Slice(wallets, func(i, j int) bool {
		return wallets[i].Id < wallets[j].Id
	})
	for i := 0; i < len(wallets); i++ {
		if wallets[i].Id != kmdWallets.Wallets[i].ID {
			t.Errorf("Wallet IDs do not match: %v ! %v", wallets[i].Id, kmdWallets.Wallets[i].ID)
		}
		if wallets[i].Name != kmdWallets.Wallets[i].Name {
			t.Errorf("Wallet names do not match: %v ! %v", wallets[i].Name, kmdWallets.Wallets[i].Name)
		}
	}
}

func TestCreate(t *testing.T) {
	kmdClient := localnet.KmdClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	t.Run("create new wallet", func(t *testing.T) {
		// Verify that the wallet does not exist
		containsWallet, err := walletManager.Contains(name)
		if err != nil {
			t.Error("Failed to list wallets", err)
			return
		}
		if containsWallet {
			t.Error("Wallet should not exist")
		}

		wallet, err := walletManager.Create(name, password)
		if err != nil {
			t.Error("Failed to create new wallet", err)
			return
		}
		if wallet.Name != name {
			t.Error("Wallet name does not match")
		}

		// Verify that the wallet was created
		containsWallet, err = walletManager.Contains(name)
		if err != nil {
			t.Error("Failed to list wallets", err)
			return
		}
		if !containsWallet {
			t.Errorf("Wallet does not exist: %#v", wallet)
		}
	})

	t.Run("create wallet with name that already exists", func(t *testing.T) {
		wallet, err := walletManager.Create(name, password)
		if err == nil {
			t.Error("Wallet creation should have failed")
		}
		t.Log(wallet, err)
		if err.Error() != "wallet with same name already exists" {
			t.Error("error message did not match")
		}
		emptyWallet := kmd.Wallet{}
		if wallet != emptyWallet {
			t.Error("wallet should have no fields set")
		}
	})

	t.Run("create wallet using blank name", func(t *testing.T) {
		wallet, err := walletManager.Create(" ", password)
		if err == nil {
			t.Error("Wallet creation should have failed")
		}
		t.Log(wallet, err)
		if err.Error() != "name cannot be blank" {
			t.Error("error message did not match")
		}
	})

	t.Run("create wallet using blank password", func(t *testing.T) {
		wallet, err := walletManager.Create(name, " ")
		if err == nil {
			t.Error("Wallet creation should have failed")
		}
		t.Log(wallet, err)
		if err.Error() != "password cannot be blank" {
			t.Error("error message did not match")
		}
	})

	t.Run("create wallet with whitespace padded name", func(t *testing.T) {
		name := ulid.Make().String()
		password := ulid.Make().String()

		wallet, err := walletManager.Create("  "+name+"  ", password)
		if err != nil {
			t.Error("Failed to create new wallet", err)
			return
		}
		if wallet.Name != name {
			t.Error("Wallet name does not match")
		}
		wallets, err := walletManager.List()
		if err != nil {
			t.Error("Failed to list wallets", err)
			return
		}
		func() {
			for _, w := range wallets {
				if w == wallet {
					return
				}
			}
			t.Errorf("Wallet does not exist: %#v", wallet)
		}()
	})

	t.Run("create wallet with whitespace padded password", func(t *testing.T) {
		name := ulid.Make().String()
		password := ulid.Make().String()

		wallet, err := walletManager.Create(name, "  "+password+"  ")
		if err != nil {
			t.Error("Failed to create new wallet", err)
			return
		}
		if wallet.Name != name {
			t.Error("Wallet name does not match")
		}
		wallets, err := walletManager.List()
		if err != nil {
			t.Error("Failed to list wallets", err)
			return
		}
		func() {
			for _, w := range wallets {
				if w == wallet {
					return
				}
			}
			t.Errorf("Wallet does not exist: %#v", wallet)
		}()
	})
}

func TestExportBackupPhrase(t *testing.T) {
	kmdClient := localnet.KmdClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	t.Run("export backup phrase for existing wallet", func(t *testing.T) {
		_, err := walletManager.Create(name, password)
		if err != nil {
			t.Error("Failed to create new wallet", err)
			return
		}
		backupPhrase, err := walletManager.ExportBackupPhrase(name, password)
		if err != nil {
			t.Error("Failed to export backup phrase", err)
			return
		}
		t.Log("backup phrase: ", backupPhrase)
	})

	t.Run("export backup phrase for wallet that does not exist", func(t *testing.T) {
		name := ulid.Make().String()
		_, err := walletManager.ExportBackupPhrase(name, password)
		if err == nil {
			t.Error("Test should have failed because the wallet does not exist")
		}
		t.Log(err)
		if err.Error() != "wallet does not exist" {
			t.Error("invalid error message")
		}
	})

	t.Run("export backup phrase using invalid password", func(t *testing.T) {
		password := ulid.Make().String()
		_, err := walletManager.ExportBackupPhrase(name, password)
		if err == nil {
			t.Error("Test should have failed because the wallet password is invalid")
		}
		t.Log(err)

		if !strings.Contains(err.Error(), "wrong password") {
			t.Error("invalid error message")
		}
	})
}

func TestRecover(t *testing.T) {
	kmdClient := localnet.KmdClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	wallet1, err := walletManager.Create(name, password)
	if err != nil {
		t.Error("Failed to create new wallet", err)
		return
	}
	backupPhrase, err := walletManager.ExportBackupPhrase(name, password)
	if err != nil {
		t.Error("Failed to export backup phrase", err)
		return
	}

	wallet2Name := name + "-2"
	wallet2, err := walletManager.Recover(wallet2Name, password, backupPhrase)
	if err != nil {
		t.Error("Failed to recover wallet", err)
		return
	}
	if wallet2.Id == wallet1.Id {
		t.Error("Recovered wallet should have a unique ID")
	}
	if wallet2.Name != wallet2Name {
		t.Error("Wallet name does not match")
	}
}

func TestCreateAccount(t *testing.T) {
	kmdClient := localnet.KmdClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	_, err := walletManager.Create(name, password)
	if err != nil {
		t.Error("Failed to create new wallet", err)
		return
	}
	accounts, err := walletManager.ListAccounts(name, password)
	if err != nil {
		t.Error("Failed to list wallet accounts", err)
		return
	}
	if len(accounts) > 0 {
		t.Error("Wallet should have no accounts")
	}

	account, err := walletManager.CreateAccount(name, password)
	if err != nil {
		t.Error("Failed to create wallet account", err)
		return
	}
	t.Log(account)

	accounts, err = walletManager.ListAccounts(name, password)
	if err != nil {
		t.Error("Failed to list wallet accounts", err)
		return
	}
	if len(accounts) != 1 {
		t.Error("Wallet should have 1 account")
	}
	if accounts[0] != account {
		t.Error("Account does not match")
	}
}

func TestCreateAccounts(t *testing.T) {
	kmdClient := localnet.KmdClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	_, err := walletManager.Create(name, password)
	if err != nil {
		t.Error("Failed to create new wallet", err)
		return
	}

	_, err = walletManager.CreateAccounts(name, password, 0)
	if err == nil {
		t.Error("Expected an error because count=0 was passed in")
	}

	var count uint = 3
	addresses, err := walletManager.CreateAccounts(name, password, count)
	if err != nil {
		t.Error("Failed to create accounts", err)
		return
	}
	t.Log(addresses)
	if len(addresses) != int(count) {
		t.Errorf("Expected number of accounts was %v but the actual number was %v", count, len(addresses))
	}
}

func TestDeleteAccount(t *testing.T) {
	kmdClient := localnet.KmdClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	_, err := walletManager.Create(name, password)
	if err != nil {
		t.Error("Failed to create new wallet", err)
		return
	}

	t.Run("delete account that exists", func(t *testing.T) {
		address, err := walletManager.CreateAccount(name, password)
		if err != nil {
			t.Error("Failed to create account", err)
			return
		}

		accountExists, err := walletManager.ContainsAccount(name, password, address)
		if err != nil {
			t.Error("Error occurred while checking if account exists", err)
			return
		}
		if !accountExists {
			t.Error("Account should exist")
		}

		if err := walletManager.DeleteAccount(name, password, address); err != nil {
			t.Error("Failed to delete account")
			return
		}
		accountExists, err = walletManager.ContainsAccount(name, password, address)
		if err != nil {
			t.Error("Failed to check for account existence", err)
			return
		}
		if accountExists {
			t.Error("Account should not exist")
		}
	})

	t.Run("delete account that does not exist", func(t *testing.T) {
		account := crypto.GenerateAccount()
		err := walletManager.DeleteAccount(name, password, account.Address.String())
		if err != nil {
			t.Error("delete failed", err)
		}

	})

}

func TestDeleteAccounts(t *testing.T) {
	kmdClient := localnet.KmdClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	_, err := walletManager.Create(name, password)
	if err != nil {
		t.Error("Failed to create new wallet", err)
		return
	}

	t.Run("Delete accounts that exist", func(t *testing.T) {
		const count = 2
		accounts := make([]string, count)
		for i := 0; i < count; i++ {
			address, err := walletManager.CreateAccount(name, password)
			if err != nil {
				t.Error("Failed to create account", err)
				return
			}
			accounts[i] = address
		}
		deleteErrors, err := walletManager.DeleteAccounts(name, password, accounts...)
		if err != nil {
			t.Error("Failed to delete accounts", err)
			return
		}
		if len(deleteErrors) > 0 {
			t.Error("There should no errors")
		}
		for _, account := range accounts {
			exists, err := walletManager.ContainsAccount(name, password, account)
			if err != nil {
				t.Error("Failed to check for account existence", err)
				return
			}
			if exists {
				t.Error("account was not deleted", account)
			}
		}

	})

	t.Run("Delete accounts that exist and don't exist", func(t *testing.T) {
		const count = 2
		accounts := make([]string, count)
		address, err := walletManager.CreateAccount(name, password)
		if err != nil {
			t.Error("Failed to create account", err)
		}
		accounts[0] = address                                   // existent account
		accounts[1] = crypto.GenerateAccount().Address.String() // non-existent account
		deleteErrors, err := walletManager.DeleteAccounts(name, password, accounts...)
		if err != nil {
			t.Error("Failed to delete accounts", err)
			return
		}
		if len(deleteErrors) > 0 {
			t.Error("There should no errors")
		}
		for _, account := range accounts {
			exists, err := walletManager.ContainsAccount(name, password, account)
			if err != nil {
				t.Error("Failed to check for account existence", err)
				return
			}
			if exists {
				t.Error("account was not deleted", account)
			}
		}

	})

	t.Run("Delete accounts that do not exist", func(t *testing.T) {
		const count = 3
		accounts := make([]string, count)
		for i := 0; i < count; i++ {
			accounts[i] = crypto.GenerateAccount().Address.String()
		}

		deleteErrors, err := walletManager.DeleteAccounts(name, password, accounts...)
		if err != nil {
			t.Error("Failed to delete accounts", err)
			return
		}
		if len(deleteErrors) > 0 {
			t.Error("There should no errors")
		}
	})
}

func TestExportPrivateKey(t *testing.T) {
	kmdClient := localnet.KmdClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	_, err := walletManager.Create(name, password)
	if err != nil {
		t.Error("Failed to create new wallet", err)
		return
	}

	t.Run("export key for account that exists", func(t *testing.T) {
		address, err := walletManager.CreateAccount(name, password)
		if err != nil {
			t.Error("Failed to create account", err)
			return
		}
		accountPassPhrase, err := walletManager.ExportPrivateKey(name, password, address)
		if err != nil {
			t.Error("Failed to export account private key", err)
			return
		}
		t.Log(accountPassPhrase)
		passPhraseWords := strings.Split(accountPassPhrase, " ")
		if len(passPhraseWords) != 25 {
			t.Error("account pass phrase mnemonic should consist of 25 words")
		}
	})

	t.Run("export key for account that does not exist", func(t *testing.T) {
		address := crypto.GenerateAccount().Address.String()
		_, err := walletManager.ExportPrivateKey(name, password, address)
		if err == nil {
			t.Error("Exporting the private key for an account that does not exist should fail")
			return
		}
		t.Log(err)
	})
}
