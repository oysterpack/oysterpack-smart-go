package core

import (
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"time"
)

type ULIDFactory interface {

	// NewULID generates a new [ULID]
	//
	// [ULID] = https://github.com/ulid/spec
	NewULID() ulid.ULID
}

func NewUlidFactory(logger *zap.Logger) ULIDFactory {
	return &cryptoULIDFactory{
		logger: logger,
	}
}

// cryptoULIDFactory implements the ULIDFactory interface
//
// It generates new ULIDs using a cryptographically secure source of entropy.
// If ULID generation fails, then the error will be logged and ULID generation will fall back to ulid.Make().
type cryptoULIDFactory struct {
	logger *zap.Logger
}

func (f *cryptoULIDFactory) NewULID() ulid.ULID {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		f.logger.Error("failed to generate ULID", zap.Error(err))
		return ulid.Make()
	}
	return id
}
