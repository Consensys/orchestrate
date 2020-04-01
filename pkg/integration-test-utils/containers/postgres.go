package containers

import (
	"context"
	"io"
	"os"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test-utils/utils"
)

type PostgresContainer struct {
	dbContainer *container.ContainerCreateCreatedBody
}

const imageName = "postgres:10.12-alpine"

func NewPostgresContainer(cli *client.Client, component string) *PostgresContainer {
	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}

	config := &container.Config{
		Image: imageName,
		Env:   []string{"POSTGRES_PASSWORD=postgres"},
		ExposedPorts: nat.PortSet{
			"5432/tcp": struct{}{},
		},
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"5432/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "5432",
				},
			},
		},
	}
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, component+"-postgres")
	if err != nil {
		panic(err)
	}

	return &PostgresContainer{dbContainer: &resp}
}

func (postgresContainer *PostgresContainer) Start(ctx context.Context, cli *client.Client) {
	utils.StartContainer(ctx, cli, postgresContainer.dbContainer.ID, 10*time.Second)
	postgres.InitEnvs()
}

func (postgresContainer *PostgresContainer) Remove(ctx context.Context, cli *client.Client) {
	utils.RemoveContainer(ctx, cli, postgresContainer.dbContainer.ID)
}

func (postgresContainer *PostgresContainer) GetDB() *pg.DB {
	return postgres.New(postgres.NewOptions())
}
