// Package app provides support for building composable applications leveraging Uber's [Fx] framework
//
// [Fx] = https://uber-go.github.io/fx/
package fxapp

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// New constructs a new Fx app instance.
//
// Params
// - options - Fx app options to apply
//
// One of the options must provide a *zap.Logger. Application logging will be configured as follows:
// - Same logger will be used for logging Fx events
// - Shutdown hook is registered to flush the logger
// - Output from the standard library's package-global logger is redirected to the supplied logger at InfoLevel
// - Global zap loggers are replaced with the provided logger
//
// NOTE: The reason the logger is not explicitly specified as a param is to allow the logger to be constructed using
// configuration and resources that ore provided by the application.

// CIa8rnHML0zaz7j
func New(options ...fx.Option) *fx.App {
	return fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Invoke(
			registerLoggerShutdownHook,
			zap.RedirectStdLog,
			zap.ReplaceGlobals,
		),
		fx.Options(options...),
	)
}

func registerLoggerShutdownHook(lc fx.Lifecycle, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = log.Sync() // ignore any errors
			return nil
		},
	})
}
