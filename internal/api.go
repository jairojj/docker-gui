package internal

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

type Api struct {
	runtime *wails.Runtime
	cli     *client.Client
	logger  *logger.CustomLogger
}

func (api *Api) WailsInit(runtime *wails.Runtime) error {
	api.runtime = runtime
	api.logger = runtime.Log.New("API")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	api.cli = cli

	go func() {
		for {
			runtime.Events.Emit("containerUpdate", api.GetContainers())
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}

func (api *Api) ListImages() []types.ImageSummary {
	ctx := context.Background()
	images, err := api.cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		api.logger.Error(err.Error())
	}

	return images
}

func (api *Api) GetContainers() []types.Container {
	ctx := context.Background()
	containers, err := api.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		api.logger.Error(err.Error())
		return nil
	}

	return containers
}
