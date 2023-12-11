package account

import (
	"context"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/oysterpack/oysterpack-smart-go/crypto/algorand/test/localnet"
	"testing"
)

func TestGetAuthAddr(t *testing.T) {
	account := localnet.GenerateTestAccount(t, types.ToMicroAlgos(0.1))
	algodClient := localnet.AlgodClient(t)
	authAddr, err := GetAuthAddr(algodClient, Address(account.Address.String()))
	if err != nil {
		t.Fatal(err)
	}
	if authAddr != Address(account.Address.String()) {
		t.Error("auth address for a new account should be the account's address")
	}
}

func TestMakeRekeyTransaction(t *testing.T) {
	account := localnet.GenerateTestAccount(t, types.ToMicroAlgos(0.2))
	authAccount := crypto.GenerateAccount()
	t.Logf("Rekey %v -> %v", account.Address, authAccount.Address)

	algodClient := localnet.AlgodClient(t)
	rekeyTxn, err := MakeRekeyTransaction(
		algodClient,
		Address(account.Address.String()),
		Address(authAccount.Address.String()),
	)
	if err != nil {
		t.Fatalf("failed to make rekey transaction: %v", err)
	}
	signer := transaction.BasicAccountTransactionSigner{
		Account: account,
	}
	signedTxns, err := signer.SignTransactions([]types.Transaction{rekeyTxn}, []int{0})
	if err != nil {
		t.Fatal(err)
	}
	txID, err := algodClient.SendRawTransaction(signedTxns[0]).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = transaction.WaitForConfirmation(algodClient, txID, 4, context.Background()); err != nil {
		t.Fatal(err)
	}

	authAddr, err := GetAuthAddr(algodClient, Address(account.Address.String()))
	if err != nil {
		t.Fatal("failed to get auth address", err)
	}
	if authAddr != Address(authAccount.Address.String()) {
		t.Errorf("auth address does not match: expected = %v, actual = %v", authAccount.Address, authAddr)
	}
}
