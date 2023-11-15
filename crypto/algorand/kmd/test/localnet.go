package test

import (
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"strings"
	"testing"
)

// LocalnetKMDClient connects to KMD localnet instance.
//
// If there is creating the client, then the test will be failed immediately.
//
// An easy way to run a localnet instance is to use [AlgoKit]
//
// [AlgoKit]: https://developer.algorand.org/docs/get-details/algokit/
func LocalnetKMDClient(t *testing.T) kmd.Client {
	const url = "http://localhost:4002"
	var token = strings.Repeat("a", 64)

	kmdClient, err := kmd.MakeClient(
		url,
		token,
	)

	if err != nil {
		t.Log("Failed to connect to KMD on localnet", err)
		t.FailNow()
	}
	return kmdClient
}
