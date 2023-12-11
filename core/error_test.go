package core

import (
	"errors"
	"github.com/oklog/ulid/v2"
	"testing"
)

var FooErrId = ulid.MustParseStrict("01HGAYQHT9S64TH85W8PXE0J02")
var ErrFoo = errors.New("foo error")

var BarErrId = ulid.MustParseStrict("01HGB2DZNSMWM5J3AV40WQ6JRW")
var ErrBar = errors.New("bar error")

var BazErrId = ulid.MustParseStrict("01HGB3N53PS74WGR2YXVXMCCKB")
var ErrBaz = errors.New("baz error")

func NewFooErr() Error {
	return Error{
		ID:   FooErrId,
		Name: "ErrFoo",
		Err:  ErrFoo,
	}

}

func NewBarErr() Error {
	return Error{
		ID:   BarErrId,
		Name: "ErrBar",
		Err:  ErrBar,
	}
}

func NewBazErr() Error {
	return Error{
		ID:   BazErrId,
		Name: "ErrBaz",
		Err:  ErrBaz,
	}
}

func TestError_Error(t *testing.T) {
	err := NewFooErr()
	t.Log(err)
	if err.ID != FooErrId {
		t.Fatal("invalid error ID")
	}

	if !errors.Is(err, NewFooErr()) {
		t.Error("errors should always match itself")
	}

	if errors.Is(err, NewBarErr()) {
		t.Error("errors should not match")
	}

	// set Error cause
	err.Cause = NewBarErr()
	t.Log(err)

	if !errors.Is(err, NewBarErr()) {
		t.Error("error should match against any error in its chain")
	}

	err.Cause = ErrBar
	if !errors.Is(err, ErrBar) {
		t.Error("error should match against any error in its chain")
	}

	// create error chain: FooErr -> BarErr -> BazErr
	cause := NewBarErr()
	cause.Cause = NewBazErr()
	err.Cause = cause
	t.Log(err)
	if !errors.Is(err, NewBazErr()) {
		t.Error("error should match against any error in its chain")
	}
}
