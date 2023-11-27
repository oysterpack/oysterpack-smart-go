package core

import (
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"testing"
)

func TestNewULID(t *testing.T) {
	// init
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal("Failed to create logger")
	}
	defer func() {
		_ = logger.Sync()
	}()
	ulidFactory := NewUlidFactory(logger)

	// Generate 10000 ULIDs and ensure they are all unique
	ulids := make(map[ulid.ULID]bool)
	const count = 10000
	var prevTime uint64 = 0
	for i := 0; i < count; i++ {
		id := ulidFactory.NewULID()
		_, ok := ulids[id]
		if ok {
			t.Fatal("duplicate ULID was generated at iteration ", i+1)
		}
		ulids[id] = true

		// every new ULID's time component should be greater than the prior one unless they were generated within the
		// same msec
		if id.Time() < prevTime {
			t.Fatal("ULID time component is older than the prior ULID")
		}
		prevTime = id.Time()
	}
	if len(ulids) != count {
		t.Fatal("The number of generated ULIDs does not match the expected value", len(ulids), "!=", count)
	}
}
