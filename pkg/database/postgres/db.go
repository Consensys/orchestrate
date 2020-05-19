package postgres

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

type DB interface {
	orm.DB
	Begin() (*pg.Tx, error)
}

type Tx interface {
	DB
	Commit() error
	Rollback() error
	Close() error
}
