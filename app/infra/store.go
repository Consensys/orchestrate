package infra

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/store/pg"
)

// InitStore initilialize envelope store
func InitStore(i *Infra) {
	opts := store.NewPGOptions()
	i.store = pg.NewEnvelopeStoreFromPGOptions(opts)
	log.WithFields(log.Fields{
		"db.address":  opts.Addr,
		"db.database": opts.Database,
		"db.user":     opts.User,
	}).Infof("infra-store: postgres store ready")
}
