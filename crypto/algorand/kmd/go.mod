module github.com/oysterpack/oysterpack-smart-go/algorand/kmd

go 1.21.4

require github.com/algorand/go-algorand-sdk/v2 v2.3.0

require github.com/oysterpack/oysterpack-smart-go/algorand/model v0.0.0-unpublished

require (
	github.com/algorand/avm-abi v0.2.0 // indirect
	github.com/algorand/go-codec/codec v1.1.10 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	golang.org/x/crypto v0.15.0 // indirect
)

replace github.com/oysterpack/oysterpack-smart-go/algorand/model => ../model
