package postgres

import (
	"context"
	"sync"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
	tlog "github.com/traefik/traefik/v2/pkg/log"
)

//go:generate mockgen -source=manager.go -destination=mocks/manager.go -package=mocks

const component = "database.postgres"

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
	mux    *sync.Mutex
	cache  map[*pg.Options]*pg.DB
	logger *log.Logger
}

func newManager() *manager {
	return &manager{
		mux:    &sync.Mutex{},
		cache:  make(map[*pg.Options]*pg.DB),
		logger: log.NewLogger().SetComponent(component),
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

	logCtx := tlog.With(
		ctx,
		tlog.Str("database", conf.Database),
		tlog.Str("addr", conf.Addr),
	)

	db := New(conf).WithContext(logCtx)

	m.logger.WithContext(logCtx).WithFields(logrus.Fields{
		"user":                 conf.User,
		"pool.size":            conf.PoolSize,
		"pool.timeout":         conf.PoolTimeout,
		"dial.timeout":         conf.DialTimeout,
		"idle.timeout":         conf.IdleTimeout,
		"idle.check-frequency": conf.IdleCheckFrequency,
	}).Info("creating database connector")

	return db
}

func (m *manager) onConnect(conn *pg.Conn) error {
	m.logger.WithContext(conn.Context()).Debug("open new connection")
	return nil
}

func GetManager() Manager {
	return mngr
}
