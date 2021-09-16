package quorumkeymanager

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	httputils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	"github.com/docker/docker/api/types/mount"
	log "github.com/sirupsen/logrus"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const defaultHashicorpVaultImage = "docker.consensys.net/pub/quorum-key-manager:v21.7.0-alpha.5"
const defaultHostPort = "8080"
const defaultHost = "localhost"

type QuorumKeyManager struct{}

type Config struct {
	Image             string
	Port              string
	MetricsPort       string
	Host              string
	DBPort            string
	DBHost            string
	ManifestDirectory string
}

func NewDefault() *Config {
	return &Config{
		Image:  defaultHashicorpVaultImage,
		Port:   defaultHostPort,
		Host:   defaultHost,
		DBHost: defaultDBHost,
		DBPort: defaultDBPort,
	}
}

func (cfg *Config) SetHostPort(port string) *Config {
	cfg.Port = port
	metricsPort, _ := strconv.ParseInt(port, 10, 64)
	cfg.MetricsPort = fmt.Sprintf("%d", metricsPort+1)
	return cfg
}

func (cfg *Config) SetHost(host string) *Config {
	if host != "" {
		cfg.Host = host
	}

	return cfg
}

func (cfg *Config) SetDBPort(port string) *Config {
	cfg.DBPort = port
	return cfg
}

func (cfg *Config) SetDBHost(host string) *Config {
	cfg.DBHost = host
	return cfg
}

func (cfg *Config) SetManifestDirectory(dir string) *Config {
	cfg.ManifestDirectory = dir
	return cfg
}

func (cfg *Config) CreateManifest(filename string, mnf *qkm.Manifest) error {
	filepath := path.Join(cfg.ManifestDirectory, filename)
	out, err := os.Create(path.Join(cfg.ManifestDirectory, filename))
	if err != nil {
		return err
	}
	defer out.Close()

	err = os.Chmod(filepath, 0777)
	if err != nil {
		return err
	}

	mnfBody, err := mnf.MarshallToYaml()
	if err != nil {
		return err
	}

	_, err = out.Write(mnfBody)
	return err
}

func (q *QuorumKeyManager) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			fmt.Sprintf("DB_PORT=%v", cfg.DBPort),
			fmt.Sprintf("DB_HOST=%v", cfg.DBHost),
		},
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
			"8081/tcp": struct{}{},
		},
		Cmd: []string{"run", "--manifest-path=/manifests", "--http-host=0.0.0.0", "--log-level=debug"},
	}

	hostConfig := &dockercontainer.HostConfig{
		RestartPolicy: dockercontainer.RestartPolicy{
			Name: "always",
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cfg.ManifestDirectory,
				Target: "/manifests",
			},
		},
	}
	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"8080/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
			"8081/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.MetricsPort}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (q *QuorumKeyManager) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*Config)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	retryT := time.NewTicker(2 * time.Second)
	defer retryT.Stop()

	httpClient := httputils.NewClient(httputils.NewDefaultConfig())

	var cerr error
waitForServiceLoop:
	for {
		select {
		case <-rctx.Done():
			cerr = rctx.Err()
			break waitForServiceLoop
		case <-retryT.C:
			resp, err := httpClient.Get(fmt.Sprintf("http://%v:%v/ready", cfg.Host, cfg.MetricsPort))
			switch {
			case err != nil:
				log.WithContext(rctx).WithError(err).Warnf("waiting for quorum-key-manager service to start")
			case resp.StatusCode != http.StatusOK:
				log.WithContext(rctx).WithField("status_code", resp.StatusCode).Warnf("waiting for quorum-key-manager service to be ready")
			default:
				log.WithContext(rctx).Info("quorum-key-manager container service is ready")
				break waitForServiceLoop
			}
		}
	}

	return cerr
}
