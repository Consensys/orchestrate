package infra

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/pg"
)

// InitStore initilialize trace store
func InitStore(i *Infra) {
	opts := infra.NewPGOptions()
	i.store = pg.NewTraceStoreFromPGOptions(opts)
	log.WithFields(log.Fields{
		"db.address":  opts.Addr,
		"db.database": opts.Database,
		"db.user":     opts.User,
	}).Infof("infra-store: postgres store ready")
}
