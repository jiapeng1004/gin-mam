package tx

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey struct{}

var txCtxKey ctxKey

func TxFromContext(ctx context.Context) (*gorm.DB, bool) {
	v := ctx.Value(txCtxKey)
	if v == nil {
		return nil, false
	}
	tx, ok := v.(*gorm.DB)
	if !ok || tx == nil {
		return nil, false
	}
	return tx, true
}

func ContextWithTx(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, txCtxKey, db)
}

func ContextWithoutTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, txCtxKey, (*gorm.DB)(nil))
}
