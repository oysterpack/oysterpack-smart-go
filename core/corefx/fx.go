package corefx

import (
	"github.com/oysterpack/oysterpack-smart-go/core"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AppMode int

const (
	Production AppMode = iota
	Development
	Testing
)

func (mode AppMode) String() string {
	switch mode {
	case Production:
		return "Production"
	case Development:
		return "Development"
	case Testing:
		return "Testing"
	default:
		return "unknown"
	}
}

func Module(mode AppMode) fx.Option {
	return fx.Module("core",
		fx.Provide(
			loggerFactory(mode),
			logger,
			core.NewUlidFactory,
		),
	)
}

// If an invalid AppMode is specified, then Production mode will be used
func loggerFactory(mode AppMode) func() core.LoggerFactory {
	switch mode {
	case Production:
		return func() core.LoggerFactory {
			return core.ProductionLoggerFactory{}
		}
	case Development:
		return func() core.LoggerFactory {
			return core.DevelopmentLoggerFactory{}
		}
	case Testing:
		return func() core.LoggerFactory {
			return core.TestingLoggerFactory{}
		}
	default:
		return func() core.LoggerFactory {
			return core.ProductionLoggerFactory{}
		}
	}
}

func logger(loggerFactory core.LoggerFactory) (*zap.Logger, error) {
	return loggerFactory.Logger()
}
