package fxapp_test

import (
	"context"
	"github.com/oysterpack/oysterpack-smart-go/fxapp"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"strings"
	"testing"

	stdlog "log"
)

var (
	onStartCounter int
	onStopCounter  int
)

func newAppLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func runApp(lc fx.Lifecycle, log *zap.Logger) {
	log.Info("Ciao mundo!")
	stdlog.Println("CIAO MUNDO!!")

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			onStartCounter++
			log.Info("APP STARTED", zap.Int("onStartCounter", onStartCounter))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			onStopCounter++
			log.Info("APP STOPPED", zap.Int("onStopCounter", onStopCounter))
			return nil
		},
	})
}

func startApp(t *testing.T, app *fx.App) {
	if err := app.Start(context.Background()); err != nil {
		t.Fatal("failed to start app", err)
	}
}

func stopApp(t *testing.T, app *fx.App) {
	if err := app.Stop(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestNew(t *testing.T) {
	t.Run("with no logger dependency", func(t *testing.T) {
		app := fxapp.New(
			fx.Invoke(runApp),
		)
		if err := app.Start(context.Background()); err == nil {
			t.Error("App should have failed to started because of missing *zap.Logger dependency")
		} else {
			if !strings.Contains(err.Error(), "missing dependencies") {
				t.Error("unexpected error:", err)
			}
		}
	})

	app := fxapp.New(
		fx.Provide(newAppLogger),
		fx.Invoke(runApp),
	)

	t.Run("start/stop app", func(t *testing.T) {
		startApp(t, app)
		stopApp(t, app)
	})

	t.Run("restart app", func(t *testing.T) {
		startApp(t, app)
		stopApp(t, app)
	})

	t.Run("starting an app that is already started should fail", func(t *testing.T) {
		startApp(t, app)
		defer stopApp(t, app)

		// try to start the app again after it has been started
		if err := app.Start(context.Background()); err == nil {
			t.Error("starting an app that is already started should fail")
		} else {
			if err.Error() != "attempted to start lifecycle when in state: started" {
				t.Error("error msg did not match:", err)
			}
		}
	})

	t.Run("stopping an app that is already stopped should be ok", func(t *testing.T) {
		startApp(t, app)
		stopApp(t, app)

		prevStopCounterValue := onStopCounter

		// try to stop the app again
		if err := app.Stop(context.Background()); err != nil {
			t.Error(err)
		}
		if onStopCounter != prevStopCounterValue {
			t.Errorf("onStopCounter should not have changed: expected: %v, actual: %v",
				prevStopCounterValue,
				onStopCounter,
			)
		}
	})

}
