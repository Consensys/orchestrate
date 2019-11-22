package hashicorp

import (
	"net"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/testutils"
)

type HashicorpKeyStoreTestSuite struct {
	testutils.SecretStoreTestSuite
}

func MockHashicorp(t *testing.T, config *Config) (net.Listener, *SecretStore) {
	t.Helper()

	// Create an in-memory, unsealed core (the "backend", if you will).
	core, _, rootToken := vault.TestCoreUnsealed(t)

	// Start an HTTP server for the core.
	ln, addr := http.TestServer(t, core)

	// Edit the provided config with the random provided address
	config.Address = addr
	// Uses default secret engine
	config.MountPoint = "secret/"
	// At the moment, kv-v2 does not work with TestCore
	config.KVVersion = "v1"

	hash, err := NewVaultClient(config)
	if err != nil {
		t.Fatal(err)
	}
	hash.Client.SetToken(rootToken)

	// Mount secret engine to use
	if config.KVVersion == "v2" {
		_ = hash.Client.Sys().TuneMount(config.MountPoint,
			api.MountConfigInput{Options: map[string]string{"version": "1"}},
		)
	}

	secretStore := &SecretStore{
		Client: hash,
		Config: config,
	}
	secretStore.ManageToken()

	return ln, secretStore
}

func SetupTest(s *HashicorpKeyStoreTestSuite, t *testing.T) {
	// TODO: Incorporate custom config
	config := ConfigFromViper()
	_, secretStore := MockHashicorp(t, config)
	s.Store = secretStore
}

func TestHashiCorp(t *testing.T) {
	s := new(HashicorpKeyStoreTestSuite)
	SetupTest(s, t)
	suite.Run(t, s)
}
