package kmd_test

import (
	"github.com/oysterpack/oysterpack-smart-go/algorand/kmd"
	"github.com/oysterpack/oysterpack-smart-go/algorand/kmd/test"
	"sort"
	"testing"
)

func TestListWallets(t *testing.T) {
	// Setup
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
	sort.Slice(kmdWallets.Wallets, func(i, j int) bool {
		return kmdWallets.Wallets[i].ID < kmdWallets.Wallets[j].ID
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
