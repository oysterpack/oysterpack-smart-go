package fxulid_test

import (
	"context"
	"github.com/oysterpack/oysterpack-smart-go/fxapp"
	"github.com/oysterpack/oysterpack-smart-go/fxulid"
	"go.uber.org/fx"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"testing"
)

func TestNewULID(t *testing.T) {
	// init
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal("Failed to create logger")
	}
	defer func() {
		_ = logger.Sync()
	}()
	newULID := fxulid.MakeNewULIDFunction(logger)

	// Generate 10000 ULIDs and ensure they are all unique
	ulids := make(map[ulid.ULID]bool)
	const count = 10000
	var prevTime uint64 = 0
	for i := 0; i < count; i++ {
		id := newULID()
		_, ok := ulids[id]
		if ok {
			t.Fatal("duplicate ULID was generated at iteration ", i+1)
		}
		ulids[id] = true

		// every new ULID's time component should be greater than the prior one unless they were generated within the
		// same msec
		if id.Time() < prevTime {
			t.Fatal("ULID time component is older than the prior ULID")
		}
		prevTime = id.Time()
	}
	if len(ulids) != count {
		t.Fatal("The number of generated ULIDs does not match the expected value", len(ulids), "!=", count)
	}
}

func newAppLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func TestMakeNewULIDFunction_WithFxApp(t *testing.T) {
	app := fxapp.New(
		fx.Provide(
			newAppLogger,
			fxulid.MakeNewULIDFunction,
		),
		fx.Invoke(func(newUlid fxulid.NewULID, log *zap.Logger) {
			log.Info("newULID()", zap.Any("ulid", newUlid()))
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	defer func(app *fx.App, ctx context.Context) {
		err := app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}(app, context.Background())
}
