package migrations

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

var contextStoreTableName = "envelopes"

func createContextTable(db migrations.DB) error {
	log.Debugf("Creating table %q...", contextStoreTableName)
	_, err := db.Exec(`CREATE TABLE ? ();`, pg.Q(contextStoreTableName))
	if err != nil {
		log.WithError(err).Errorf("Could not create table %q", contextStoreTableName)
		return err
	}
	log.Infof("Created table %q", contextStoreTableName)

	return nil
}

func dropContextTable(db migrations.DB) error {
	log.Debugf("Dropping table %q...", contextStoreTableName)
	_, err := db.Exec(`DROP TABLE ?;`, pg.Q(contextStoreTableName))
	if err != nil {
		log.WithError(err).Errorf("Could not drop table %q", contextStoreTableName)
		return err
	}
	log.Infof("Dropped table %q", contextStoreTableName)

	return nil
}

func init() {
	Collection.MustRegisterTx(createContextTable, dropContextTable)
}
