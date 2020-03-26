package postgres

import (
	"context"

	"github.com/go-pg/pg/v9"
)

type pgCtxKey string

const txKey pgCtxKey = "tx"

func WithTx(ctx context.Context, tx *pg.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func TxFromContext(ctx context.Context) *pg.Tx {
	tx, _ := ctx.Value(txKey).(*pg.Tx)
	return tx
}
