package kafka

import (
	"context"
	"fmt"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const DefaultKafkaImage = "confluentinc/cp-kafka:5.3.0"

type Kafka struct{}
type Config struct {
	Image               string
	Port                string
	ZookeeperClientPort string
	ZookeeperHostname   string
	KafkaHostname       string
}

func (p *Config) SetDefault() *Config {
	if p.Image == "" {
		p.Image = DefaultKafkaImage
	}

	if p.Port == "" {
		p.Port = "9092"
	}

	if p.ZookeeperClientPort == "" {
		p.ZookeeperClientPort = "32181"
	}

	if p.ZookeeperHostname == "" {
		p.ZookeeperHostname = "zookeeper"
	}

	if p.KafkaHostname == "" {
		p.KafkaHostname = "kafka"
	}

	return p
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
			fmt.Sprintf("KAFKA_ADVERTISED_LISTENERS=INTERNAL://%v:29092,EXTERNAL://localhost:%v", cfg.KafkaHostname, cfg.Port),
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT",
			"KAFKA_INTER_BROKER_LISTENER_NAME=INTERNAL",
		},
		ExposedPorts: nat.PortSet{
			"9092/tcp": struct{}{},
		},
	}

	hostConfig := &dockercontainer.HostConfig{
		PortBindings: nat.PortMap{
			"9092/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		},
	}

	return containerCfg, hostConfig, nil, nil
}
