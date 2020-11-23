package keymanager

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	// Initialize dependencies
	config := NewConfig(viper.GetViper())

	return NewKeyManager(ctx, config)
}

func Run(ctx context.Context) error {
	appli, err := New(ctx)
	if err != nil {
		return err
	}
	return appli.Run(ctx)
}
