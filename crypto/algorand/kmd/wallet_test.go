package kmd_test

import (
	"github.com/oklog/ulid/v2"
	"github.com/oysterpack/oysterpack-smart-go/crypto/algorand/kmd"
	"github.com/oysterpack/oysterpack-smart-go/crypto/algorand/kmd/test"
	"sort"
	"strings"
	"testing"
)

func TestListWallets(t *testing.T) {
	kmdClient := test.LocalnetKMDClient(t)

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
	kmdClient := test.LocalnetKMDClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	t.Run("create new wallet", func(t *testing.T) {
		wallet, err := walletManager.Create(name, password)
		if err != nil {
			t.Error("Failed to create new wallet", err)
		}
		if wallet.Name != name {
			t.Error("Wallet name does not match")
		}
		wallets, err := walletManager.List()
		if err != nil {
			t.Error("Failed to list wallets", err)
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
		}
		if wallet.Name != name {
			t.Error("Wallet name does not match")
		}
		wallets, err := walletManager.List()
		if err != nil {
			t.Error("Failed to list wallets", err)
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
		}
		if wallet.Name != name {
			t.Error("Wallet name does not match")
		}
		wallets, err := walletManager.List()
		if err != nil {
			t.Error("Failed to list wallets", err)
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
	kmdClient := test.LocalnetKMDClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	t.Run("export backup phrase for existing wallet", func(t *testing.T) {
		_, err := walletManager.Create(name, password)
		if err != nil {
			t.Error("Failed to create new wallet", err)
		}
		backupPhrase, err := walletManager.ExportBackupPhrase(name, password)
		if err != nil {
			t.Error("Failed to export backup phrase", err)
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
	kmdClient := test.LocalnetKMDClient(t)
	walletManager := kmd.New(kmdClient)

	name := ulid.Make().String()
	password := ulid.Make().String()

	wallet1, err := walletManager.Create(name, password)
	if err != nil {
		t.Error("Failed to create new wallet", err)
	}
	backupPhrase, err := walletManager.ExportBackupPhrase(name, password)
	if err != nil {
		t.Error("Failed to export backup phrase", err)
	}

	wallet2Name := name + "-2"
	wallet2, err := walletManager.Recover(wallet2Name, password, backupPhrase)
	if err != nil {
		t.Error("Failed to recover wallet", err)
	}
	if wallet2.Id == wallet1.Id {
		t.Error("Recovered wallet should have a unique ID")
	}
	if wallet2.Name != wallet2Name {
		t.Error("Wallet name does not match")
	}
}
