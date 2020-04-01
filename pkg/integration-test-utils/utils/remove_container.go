package utils

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RemoveContainer(ctx context.Context, cli *client.Client, id string) {
	if err := cli.ContainerStop(ctx, id, nil); err != nil {
		panic(err)
	}

	if err := cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
	}); err != nil {
		panic(err)
	}

	fmt.Println("Container removed successfully, ID:", id)
}
