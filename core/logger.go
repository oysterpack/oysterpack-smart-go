package core

import (
	"go.uber.org/zap"
)

type LoggerFactory interface {
	Logger() (*zap.Logger, error)
}

type ProductionLoggerFactory struct{}

func (l ProductionLoggerFactory) Logger() (*zap.Logger, error) {
	return zap.NewProduction()
}

type DevelopmentLoggerFactory struct{}

func (l DevelopmentLoggerFactory) Logger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

type TestingLoggerFactory struct{}

func (l TestingLoggerFactory) Logger() (*zap.Logger, error) {
	return zap.NewExample(), nil
}
