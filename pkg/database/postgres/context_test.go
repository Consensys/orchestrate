package postgres

import (
	"context"
	"testing"

	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := WithTx(context.Background(), &pg.Tx{})
	tx := TxFromContext(ctx)
	assert.Equal(t, &pg.Tx{}, tx)
}
