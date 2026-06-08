package tx

import (
	"context"

	"gorm.io/gorm"
)

//go:generate mockgen -source=manager.go -destination=mock/manager_mock.go -package=mock

type Manager interface {
	Run(ctx context.Context, prop Propagation, fn func(ctx context.Context) error) error
}

type manager struct {
	db *gorm.DB
}

func NewManager(db *gorm.DB) Manager {
	return &manager{db: db}
}

func (m *manager) Run(ctx context.Context, prop Propagation, fn func(ctx context.Context) error) error {
	switch prop {
	case Required:
		if _, ok := TxFromContext(ctx); ok {
			return fn(ctx)
		}
		return m.db.Transaction(func(tx *gorm.DB) error {
			return fn(ContextWithTx(ctx, tx))
		})
	case RequiresNew:
		return m.db.Transaction(func(tx *gorm.DB) error {
			return fn(ContextWithTx(ctx, tx))
		})
	case Mandatory:
		if _, ok := TxFromContext(ctx); !ok {
			return ErrTxMandatory
		}
		return fn(ctx)
	case NotSupported:
		return fn(ContextWithoutTx(ctx))
	case Never:
		if _, ok := TxFromContext(ctx); ok {
			return ErrTxNever
		}
		return fn(ctx)
	default:
		return ErrTxUnsupportedPropagation
	}
}
