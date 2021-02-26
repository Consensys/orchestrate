package manager

import (
	"context"
	"crypto/tls"
)

//go:generate mockgen -source=manager.go -destination=mock/mock.go -package=mock

type Manager interface {
	Get(ctx context.Context, configuration interface{}) (*tls.Config, map[string]*tls.Config, error)
}
