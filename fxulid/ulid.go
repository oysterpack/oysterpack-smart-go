package fxulid

import (
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"time"
)

// NewULID generates new [ULID]s
//
// [ULID] = https://github.com/ulid/spec
type NewULID func() ulid.ULID

// MakeNewULIDFunction is a NewULID constructor function
//
// The returned NewULID function generates new ULIDs using a cryptographically secure source of entropy.
// If ULID generation fails, then the error will be logged and ULID generation will fall back to ulid.Make().
func MakeNewULIDFunction(logger *zap.Logger) NewULID {
	return func() ulid.ULID {
		id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
		if err != nil {
			logger.Error("failed to generate ULID", zap.Error(err))
			return ulid.Make()
		}
		return id
	}
}
