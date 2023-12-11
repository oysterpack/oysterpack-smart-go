package corefx

import (
	"context"
	"github.com/oysterpack/oysterpack-smart-go/core"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"testing"
)

func TestNewFx(t *testing.T) {
	test := func(loggerConstructor LoggerProvider) func(t *testing.T) {
		return func(t *testing.T) {
			app := NewFx(
				loggerConstructor,                 // *zap.Logger
				fx.Provide(core.NewULIDGenerator), // ULIDFactory
				fx.Invoke(func(logger *zap.Logger, newULID core.ULIDGenerator) {
					id := newULID()
					logger.Info("generated new ULID", zap.Any("ulid", id))
				}),
			)

			if err := app.Start(context.Background()); err != nil {
				t.Fatal("failed to start app", err)
			}

			if err := app.Stop(context.Background()); err != nil {
				t.Fatal("failed to stop app", err)
			}
		}
	}

	t.Run("with DevLogger", test(DevLogger))
	t.Run("with ProdLogger", test(ProdLogger))

}
