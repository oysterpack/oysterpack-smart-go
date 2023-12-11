package core

import (
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
)

type Error struct {
	ID    ulid.ULID // unique error ID
	Name  string    // human friendly name (naming convention is to prefix the name with "Err")
	Err   error     // underlying error
	Cause error     // error chain
}

func (e Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%v[%v]: %v <> %v", e.Name, e.ID, e.Err, e.Cause)
	}
	return fmt.Sprintf("%v[%v]: %v", e.Name, e.ID, e.Err)
}

func (e Error) Unwrap() error {
	return e.Cause
}

// Is reports true if the err type is Error and if the ID matches.
// Otherwise, check the error against its cause.
func (e Error) Is(err error) bool {
	switch err := err.(type) {
	case Error:
		return e.ID == err.ID
	default:
		return errors.Is(err, e.Cause)
	}
}
