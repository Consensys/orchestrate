package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/config"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/compose"
)

type Client struct {
	cli *client.Client

	composition *config.Composition
	generator   container.ConfigGenerator
	containers  map[string]dockercontainer.ContainerCreateCreatedBody
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
		generator:   compose.New(),
		containers:  make(map[string]dockercontainer.ContainerCreateCreatedBody),
	}, nil
}

func (c *Client) Up(ctx context.Context, name string) error {
	logger := log.FromContext(ctx).WithField("container", name)

	containerCfg, hostCfg, networkCfg, err := c.generator.GenerateContainerConfig(ctx, c.composition.Containers[name])
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

	logger.WithField("image", containerCfg.Image).Infof("pulled imaged")

	// Create Docker container
	containerBody, err := c.cli.ContainerCreate(ctx, containerCfg, hostCfg, networkCfg, name)
	if err != nil {
		return err
	}
	c.containers[name] = containerBody

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

func (c *Client) Down(ctx context.Context, name string) error {
	logger := log.FromContext(ctx).WithField("container", name)

	containerBody, err := c.getContainer(name)
	if err != nil {
		return nil
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

func (c *Client) getContainer(name string) (dockercontainer.ContainerCreateCreatedBody, error) {
	containerBody, ok := c.containers[name]
	if !ok {
		return dockercontainer.ContainerCreateCreatedBody{}, fmt.Errorf("no container named %v", name)
	}

	return containerBody, nil
}
