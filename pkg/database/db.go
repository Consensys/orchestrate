package database

type DB interface {
	Begin() (Tx, error)
}

type Tx interface {
	DB
	Commit() error
	Rollback() error
	Close() error
}

func ExecuteInDBTx(db DB, persist func(tx Tx) error) (err error) {
	dbtx, isTx := db.(Tx)
	if !isTx {
		if dbtx, err = db.Begin(); err != nil {
			return err
		}
	}

	defer func() {
		// In case it is a nested transaction IGNORE
		if isTx {
			return
		}

		if err != nil {
			_ = dbtx.Rollback()
		} else {
			err = dbtx.Commit()
		}
	}()

	if err := persist(dbtx); err != nil {
		return err
	}

	return nil
}
