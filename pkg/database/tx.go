package database

type Tx interface {
	Commit() error
}
