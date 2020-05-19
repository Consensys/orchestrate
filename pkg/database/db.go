package database

import (
	"reflect"
)

type DB interface {
	Begin() (Tx, error)
}

type Tx interface {
	DB
	Commit() error
	Rollback() error
	Close() error
}

func ExecuteInDBTx(db DB, persist func(tx Tx) error) (der error) {
	dbtx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		// In case it is a nested transaction IGNORE
		if reflect.DeepEqual(db, dbtx) {
			return
		}

		if der == nil {
			der = dbtx.Commit()
		}

		if der != nil {
			if err := dbtx.Rollback(); err != nil {
				der = err
			}
		}

		if err := dbtx.Close(); err != nil {
			der = err
		}
	}()

	if err := persist(dbtx); err != nil {
		return err
	}

	return nil
}
