package corefx

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// NewFx constructs a new Fx application using the specified loggerConstructor and options
// - the same Zap logger will be used for Fx logs as well
func NewFx(loggerConstructor LoggerProvider, options ...fx.Option) *fx.App {
	return fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(loggerConstructor),
		fx.Options(options...),
	)
}

type LoggerProvider func() (*zap.Logger, error)

// DevLogger is a LoggerProvider, which is meant to be used for development purposes.
// It logs at DebugLevel.
func DevLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

// ProdLogger builds a sensible production Logger that writes InfoLevel and above logs to standard error as JSON.
func ProdLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
