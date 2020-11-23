package kafka

import (
	"context"
	"fmt"
	"time"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/zookeeper"
)

const DefaultKafkaImage = "confluentinc/cp-kafka:5.3.0"
const DefaultExternalHostname = "localhost:9092"

type Kafka struct {
}

type Config struct {
	Image                 string
	Port                  string
	ZookeeperClientPort   string
	ZookeeperHostname     string
	KafkaInternalHostname string
	KafkaExternalHostname string
}

func NewDefault() *Config {
	return &Config{
		Image:                 DefaultKafkaImage,
		ZookeeperClientPort:   zookeeper.DefaultZookeeperClientPort,
		ZookeeperHostname:     "zookeeper",
		KafkaInternalHostname: "kafka",
		KafkaExternalHostname: DefaultExternalHostname,
	}
}

func (cfg *Config) SetHostPort(port string) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetZookeeperHostname(host string) *Config {
	cfg.ZookeeperHostname = host
	return cfg
}

func (cfg *Config) SetKafkaInternalHostname(name string) *Config {
	cfg.KafkaInternalHostname = name
	return cfg
}

func (cfg *Config) SetKafkaExternalHostname(name string) *Config {
	cfg.KafkaExternalHostname = name
	return cfg
}

func (k *Kafka) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			"KAFKA_BROKER_ID=1",
			fmt.Sprintf("KAFKA_ZOOKEEPER_CONNECT=%v:%v", cfg.ZookeeperHostname, cfg.ZookeeperClientPort),
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
			"KAFKA_LISTENERS=INTERNAL://0.0.0.0:29092,EXTERNAL://0.0.0.0:9092",
			fmt.Sprintf("KAFKA_ADVERTISED_LISTENERS=INTERNAL://%v:29092,EXTERNAL://%v", cfg.KafkaInternalHostname, cfg.KafkaExternalHostname),
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT",
			"KAFKA_INTER_BROKER_LISTENER_NAME=INTERNAL",
		},
		ExposedPorts: nat.PortSet{
			"9092/tcp":  struct{}{},
			"29092/tcp": struct{}{},
		},
	}

	hostConfig := &dockercontainer.HostConfig{}
	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"9092/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (k *Kafka) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*Config)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	saramaCfg, _ := pkgsarama.NewSaramaConfig()
	addrs := []string{cfg.KafkaExternalHostname}

	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	retryT := time.NewTicker(time.Second)
	defer retryT.Stop()

	var cerr error
waitForServiceLoop:
	for {
		select {
		case <-rctx.Done():
			cerr = rctx.Err()
			break waitForServiceLoop
		case <-retryT.C:
			client, err := pkgsarama.NewClient(addrs, saramaCfg)
			switch {
			case err != nil:
				log.WithContext(rctx).
					WithError(err).
					Warnf("waiting for kafka service to start: %s", cfg.KafkaExternalHostname)
			case len(client.Brokers()) < 1:
				err := fmt.Errorf("not available brokers")
				log.WithContext(rctx).
					WithError(err).
					Warnf("waiting for kafka service to start: %s", cfg.KafkaExternalHostname)
			default:
				log.WithContext(rctx).Infof("kafka container service is ready")
				break waitForServiceLoop
			}
		}
	}

	return cerr
}
