module github.com/oysterpack/oysterpack-smart-go/crypto/algorand

go 1.21.4

require (
	github.com/algorand/go-algorand-sdk/v2 v2.3.0
	github.com/oklog/ulid/v2 v2.1.0
	github.com/oysterpack/oysterpack-smart-go/core v0.0.0-unpublished
)

require (
	github.com/algorand/avm-abi v0.1.1 // indirect
	github.com/algorand/go-codec/codec v1.1.10 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/crypto v0.15.0 // indirect
)

replace github.com/oysterpack/oysterpack-smart-go/core v0.0.0-unpublished => ../../core
