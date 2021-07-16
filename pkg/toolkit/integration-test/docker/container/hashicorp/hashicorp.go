package hashicorp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types/mount"
	api2 "github.com/hashicorp/vault/api"

	httputils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	log "github.com/sirupsen/logrus"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const defaultHashicorpVaultImage = "library/vault:1.6.2"
const defaultHostPort = "8200"
const defaultRootToken = "myRoot"
const defaultHost = "localhost"
const pluginFileName = "orchestrate-hashicorp-vault-plugin"
const pluginVersion = "v0.0.11-alpha.3"
const defaultMountPath = "orchestrate"

type Vault struct{}

type Config struct {
	Image                 string
	Host                  string
	Port                  string
	RootToken             string
	MonthPath             string
	PluginSourceDirectory string
}

func NewDefault() *Config {
	return &Config{
		Image:     defaultHashicorpVaultImage,
		Port:      defaultHostPort,
		RootToken: defaultRootToken,
		Host:      defaultHost,
		MonthPath: defaultMountPath,
	}
}

func (cfg *Config) SetHostPort(port string) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetRootToken(rootToken string) *Config {
	cfg.RootToken = rootToken
	return cfg
}

func (cfg *Config) SetHost(host string) *Config {
	if host != "" {
		cfg.Host = host
	}

	return cfg
}

func (cfg *Config) SetMountPath(mountPath string) *Config {
	cfg.MonthPath = mountPath
	return cfg
}

func (cfg *Config) SetPluginSourceDirectory(dir string) *Config {
	cfg.PluginSourceDirectory = dir
	return cfg
}

func (cfg *Config) DownloadPlugin() error {
	url := fmt.Sprintf("https://github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/releases/download/%s/%s", pluginVersion, pluginFileName)
	err := downloadPlugin(fmt.Sprintf("%s/%s", cfg.PluginSourceDirectory, pluginFileName), url)
	if err != nil {
		return err
	}
	return nil
}

func (vault *Vault) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			fmt.Sprintf("VAULT_DEV_ROOT_TOKEN_ID=%v", cfg.RootToken),
		},
		ExposedPorts: nat.PortSet{
			"8200/tcp": struct{}{},
		},
		Tty: true,
		Cmd: []string{"server", "-dev", "-dev-plugin-dir=/vault/plugins", "-log-level=debug"},
	}

	hostConfig := &dockercontainer.HostConfig{
		CapAdd: []string{"IPC_LOCK"},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cfg.PluginSourceDirectory,
				Target: "/vault/plugins",
			},
		},
	}
	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"8200/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (vault *Vault) enablePlugin(serverAddr, rootToken, mountPath string) error {
	// Enable orchestrate secret engine
	vaultClient, err := api2.NewClient(&api2.Config{
		Address: serverAddr,
	})
	if err != nil {
		return err
	}

	vaultClient.SetToken(rootToken)
	return vaultClient.Sys().Mount(mountPath, &api2.MountInput{
		Type:        "plugin",
		Description: "Orchestrate Wallets",
		Config: api2.MountConfigInput{
			ForceNoCache:              true,
			PassthroughRequestHeaders: []string{"X-Vault-Namespace"},
		},
		PluginName: pluginFileName,
	})
}

func (vault *Vault) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*Config)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	retryT := time.NewTicker(2 * time.Second)
	defer retryT.Stop()

	httpClient := httputils.NewClient(httputils.NewDefaultConfig())
	serverAddr := "http://" + cfg.Host + ":" + cfg.Port

	var cerr error
waitForServiceLoop:
	for {
		select {
		case <-rctx.Done():
			cerr = rctx.Err()
			break waitForServiceLoop
		case <-retryT.C:
			resp, err := httpClient.Get(fmt.Sprintf("%s/v1/sys/health", serverAddr))

			switch {
			case err != nil:
				log.WithContext(rctx).WithError(err).Warnf("waiting for Hashicorp Vault service to start")
			case resp.StatusCode != http.StatusOK:
				log.WithContext(rctx).WithField("status_code", resp.StatusCode).Warnf("waiting for Hashicorp Vault service to be ready")
			default:
				log.WithContext(rctx).Info("hashicorp Vault container service is ready")
				break waitForServiceLoop
			}
		}
	}

	if cerr != nil {
		return cerr
	}

	if err := vault.enablePlugin(serverAddr, cfg.RootToken, cfg.MonthPath); err != nil {
		return err
	}

	return nil
}

func downloadPlugin(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = os.Chmod(filepath, 0777)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	return err
}
