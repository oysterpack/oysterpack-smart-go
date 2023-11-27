package core

import (
	"testing"
)

func TestDevelopmentLoggerFactory_Logger(t *testing.T) {
	// DevelopmentLoggerFactory should implement the LoggerFactory interface
	var loggerFactory LoggerFactory = DevelopmentLoggerFactory{}
	logger, err := loggerFactory.Logger()
	if err != nil {
		t.Fatal("failed to create logger")
	}
	defer func() {
		// ignore error
		_ = logger.Sync()
	}()
	logger.Info("dev")
}

func TestProductionLoggerFactory_Logger(t *testing.T) {
	// ProductionLoggerFactory should implement the LoggerFactory interface
	var loggerFactory LoggerFactory = ProductionLoggerFactory{}
	logger, err := loggerFactory.Logger()
	if err != nil {
		t.Fatal("failed to create logger")
	}
	defer func() {
		// ignore error
		_ = logger.Sync()
	}()
	logger.Info("prod")
}

func TestTestingLoggerFactory_Logger(t *testing.T) {
	// TestingLoggerFactory should implement the LoggerFactory interface
	var loggerFactory LoggerFactory = TestingLoggerFactory{}
	logger, err := loggerFactory.Logger()
	if err != nil {
		t.Fatal("failed to create logger")
	}
	defer func() {
		// ignore error
		_ = logger.Sync()
	}()
	logger.Info("test")
}
