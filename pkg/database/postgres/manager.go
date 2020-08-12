package postgres

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=manager.go -destination=mocks/manager.go -package=mocks

func init() {
	mngr = newManager()
}

var mngr *manager

type Manager interface {
	Connect(ctx context.Context, conf *pg.Options) *pg.DB
}

func NewManager() Manager {
	return newManager()
}

type manager struct {
	mux   *sync.Mutex
	cache map[*pg.Options]*pg.DB
}

func newManager() *manager {
	return &manager{
		mux:   &sync.Mutex{},
		cache: make(map[*pg.Options]*pg.DB),
	}
}

func (m *manager) Connect(ctx context.Context, conf *pg.Options) *pg.DB {
	m.mux.Lock()
	defer m.mux.Unlock()
	if db, ok := m.cache[conf]; ok {
		return db
	}

	if conf.OnConnect == nil {
		conf.OnConnect = m.onConnect
	}

	logCtx := log.With(
		ctx,
		log.Str("database", conf.Database),
		log.Str("addr", conf.Addr),
	)

	db := New(conf).WithContext(logCtx)

	log.FromContext(logCtx).
		WithFields(logrus.Fields{
			"user":                 conf.User,
			"pool.size":            conf.PoolSize,
			"pool.timeout":         conf.PoolTimeout,
			"dial.timeout":         conf.DialTimeout,
			"idle.timeout":         conf.IdleTimeout,
			"idle.check-frequency": conf.IdleCheckFrequency,
		}).
		Infof("creating postgres database connector")

	return db
}

func (m *manager) onConnect(conn *pg.Conn) error {
	log.FromContext(conn.Context()).Debugf("open connection to postgres")
	return nil
}

func GetManager() Manager {
	return mngr
}
