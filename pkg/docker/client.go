package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/network"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/config"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/compose"
)

type Client struct {
	cli *client.Client

	composition *config.Composition
	factory     container.DockerContainerFactory
	containers  map[string]dockercontainer.ContainerCreateCreatedBody
	networks    map[string]string
}

func NewClient(composition *config.Composition) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		cli:         cli,
		composition: composition,
		factory:     compose.New(),
		containers:  make(map[string]dockercontainer.ContainerCreateCreatedBody),
		networks:    make(map[string]string),
	}, nil
}

func (c *Client) Up(ctx context.Context, name, networkName string) error {
	logger := log.FromContext(ctx).WithField("container", name)

	containerCfg, hostCfg, networkCfg, err := c.factory.GenerateContainerConfig(ctx, c.composition.Containers[name])
	if err != nil {
		return err
	}

	// Pull image
	reader, err := c.cli.ImagePull(ctx, containerCfg.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}
	_ = reader.Close()
	logger.WithField("image", containerCfg.Image).Info("pulled imaged")

	// Create Docker container
	containerBody, err := c.cli.ContainerCreate(ctx, containerCfg, hostCfg, networkCfg, name)
	if err != nil {
		return err
	}
	c.containers[name] = containerBody

	// Connect to network and assign the alias as the name of the container
	if networkID, ok := c.networks[networkName]; networkName != "" && ok {
		err = c.cli.NetworkConnect(ctx, networkID, containerBody.ID, &network.EndpointSettings{
			Aliases:   []string{name},
			NetworkID: networkID,
		})

		if err != nil {
			return err
		}

		logger.WithField("network_id", networkID).Infof("container %v connected to network", name)
	} else if networkName != "" {
		errMsg := fmt.Sprintf("container %v cannot connected to network", name)
		logger.WithField("network_id", networkID).Errorf(errMsg)
		return errors.InvalidArgError(errMsg)
	}

	// Start docker container
	if err := c.cli.ContainerStart(ctx, containerBody.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	logger.WithField("id", containerBody.ID).Info("started container")

	return nil
}

func (c *Client) Start(ctx context.Context, name string) error {
	logger := log.FromContext(ctx).WithField("container", name)

	containerBody, err := c.getContainer(name)
	if err != nil {
		return nil
	}

	// Start docker container
	if err := c.cli.ContainerStart(ctx, containerBody.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	logger.WithField("id", containerBody.ID).Infof("started container")

	return nil
}

func (c *Client) Stop(ctx context.Context, name string) error {
	logger := log.FromContext(ctx).WithField("container", name)

	containerBody, err := c.getContainer(name)
	if err != nil {
		return nil
	}

	if err := c.cli.ContainerStop(ctx, containerBody.ID, nil); err != nil {
		return err
	}

	logger.WithField("id", containerBody.ID).Infof("stopped container")

	return nil
}

func (c *Client) WaitTillIsReady(ctx context.Context, name string, timeout time.Duration) error {
	logger := log.FromContext(ctx).WithField("container", name)

	err := c.factory.WaitForService(ctx, c.composition.Containers[name], timeout)
	if err != nil {
		logger.WithError(err).Errorf("cannot wait for service %s", name)
		return err
	}

	return nil
}

func (c *Client) StartServiceAndWait(ctx context.Context, name string, timeout time.Duration) error {
	err := c.Start(ctx, name)
	if err != nil {
		return err
	}

	err = c.WaitTillIsReady(ctx, name, timeout)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Down(ctx context.Context, name string) error {
	logger := log.FromContext(ctx).WithField("container", name)

	containerBody, err := c.getContainer(name)
	if err != nil {
		logger.WithError(err).Warnf("container %s cannot be found", name)
		if err := c.cli.ContainerStop(ctx, name, nil); err != nil {
			return err
		}
	}

	if err := c.cli.ContainerStop(ctx, containerBody.ID, nil); err != nil {
		return err
	}

	logger.WithField("id", containerBody.ID).Infof("stopped container")

	if err := c.cli.ContainerRemove(ctx, containerBody.ID, types.ContainerRemoveOptions{RemoveVolumes: true}); err != nil {
		return err
	}

	logger.WithField("id", containerBody.ID).Infof("removed container")

	return nil
}

func (c *Client) CreateNetwork(ctx context.Context, name string) error {
	logger := log.FromContext(ctx).WithField("network_name", name)

	createResponse, err := c.cli.NetworkCreate(ctx, name, types.NetworkCreate{Driver: "bridge"})
	if err != nil {
		return err
	}

	logger.WithField("id", createResponse.ID).Infof("created network")
	c.networks[name] = createResponse.ID
	return nil
}

func (c *Client) RemoveNetwork(ctx context.Context, name string) error {
	logger := log.FromContext(ctx).WithField("network_name", name)

	if networkID, ok := c.networks[name]; !ok {
		err := c.cli.NetworkRemove(ctx, networkID)
		if err != nil {
			return err
		}

		logger.WithField("network_id", networkID).Infof("removed network")
	} else {
		err := c.cli.NetworkRemove(ctx, name)
		if err != nil {
			return err
		}

		logger.WithField("network_name", name).Infof("removed network")
	}

	return nil
}

func (c *Client) getContainer(name string) (dockercontainer.ContainerCreateCreatedBody, error) {
	containerBody, ok := c.containers[name]
	if !ok {
		return dockercontainer.ContainerCreateCreatedBody{}, fmt.Errorf("no container found with name '%v'", name)
	}

	return containerBody, nil
}
