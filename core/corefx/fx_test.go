package corefx

import (
	"context"
	"github.com/oysterpack/oysterpack-smart-go/core"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"testing"
)

func TestModule(t *testing.T) {
	modes := []AppMode{
		Production,
		Development,
		Testing,
	}

	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			app := fx.New(
				fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
					return &fxevent.ZapLogger{Logger: log}
				}),
				Module(mode),
				fx.Invoke(func(logger *zap.Logger, ulidFactory core.ULIDFactory) {
					ulid := ulidFactory.NewULID()
					logger.Info("generated new ULID", zap.Any("ulid", ulid))
				}),
			)

			if err := app.Start(context.Background()); err != nil {
				t.Fatal("failed to start app", err)
			}

			if err := app.Stop(context.Background()); err != nil {
				t.Fatal("failed to stop app", err)
			}
		})

	}

}
