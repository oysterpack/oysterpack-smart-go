package account

import (
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/oysterpack/oysterpack-smart-go/core"
)

var (
	ErrGetAuthAddrFailed        = ulid.MustParse("01HGB4K37WVFBQFPM0RND212YS")
	ErrAccountAlreadyRekeyed    = ulid.MustParse("01HGB4MAW16F6GXHCKC3399BZ1")
	ErrGetSuggestedParamsFailed = ulid.MustParse("01HGTE11YD3XAX2KGZVGWZFWRY")
	ErrMakePaymentTxn           = ulid.MustParse("01HGTEPA5PWKCTXM682GRBJCVB")
	ErrSignTransactions         = ulid.MustParse("01HGTETWSHXMFNSTFWMS5JDZRZ")
	ErrSettingRekeyTo           = ulid.MustParse("01HGTF3KG43JWPT8SBCF7W67WX")
)

func errGetAuthAddrFailed(cause error) core.Error {
	return core.Error{
		ID:    ErrGetAuthAddrFailed,
		Name:  "ErrGetAuthAddrFailed",
		Err:   errors.New("failed to get account auth address"),
		Cause: cause,
	}
}

func errAccountAlreadyRekeyed(address Address) core.Error {
	return core.Error{
		ID:   ErrAccountAlreadyRekeyed,
		Name: "ErrAccountAlreadyRekeyed",
		Err:  fmt.Errorf("account has already been rekeyed: %v", address),
	}
}

func errGetSuggestedParamsFailed(cause error) core.Error {
	return core.Error{
		ID:    ErrGetSuggestedParamsFailed,
		Name:  "ErrGetSuggestedParamsFailed",
		Err:   errors.New("failed to get suggested params for constructing a new transaction"),
		Cause: cause,
	}
}

func errMakePaymentTxn(cause error) core.Error {
	return core.Error{
		ID:    ErrMakePaymentTxn,
		Name:  "ErrMakePaymentTxn",
		Err:   errors.New("failed to construct payment transaction"),
		Cause: cause,
	}
}

func errSignTransactions(cause error) core.Error {
	return core.Error{
		ID:    ErrSignTransactions,
		Name:  "ErrSignTransactions",
		Err:   errors.New("failed to sign transactions"),
		Cause: cause,
	}
}

func errSettingRekeyTo(address Address, cause error) core.Error {
	return core.Error{
		ID:    ErrSettingRekeyTo,
		Name:  "ErrSettingRekeyTo",
		Err:   fmt.Errorf("failed to set the rekeyTo field on the transaction: %v", address),
		Cause: cause,
	}
}
