module github.com/oysterpack/oysterpack-smart-go/fxulid

go 1.21.4

require (
	github.com/oklog/ulid/v2 v2.1.0
	github.com/oysterpack/oysterpack-smart-go/fxapp v0.0.0-unpublished
	go.uber.org/fx v1.20.1
	go.uber.org/zap v1.26.0
)

require (
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
)

replace github.com/oysterpack/oysterpack-smart-go/fxapp v0.0.0-unpublished => ../fxapp
