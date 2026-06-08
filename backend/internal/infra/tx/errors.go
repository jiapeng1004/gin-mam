package tx

import "errors"

var (
	ErrTxMandatory                = errors.New("transaction mandatory but none active")
	ErrTxNever                    = errors.New("transaction must not run within an existing transaction")
	ErrTxUnsupportedPropagation = errors.New("unsupported transaction propagation")
)
