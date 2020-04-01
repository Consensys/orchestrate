package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func StartContainer(ctx context.Context, cli *client.Client, id string, sleep time.Duration) {
	fmt.Println("Container starting, ID:", id)
	if err := cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	time.Sleep(sleep)
	fmt.Println("Container started successfully, ID:", id)
}
