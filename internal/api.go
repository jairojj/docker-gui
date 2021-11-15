package internal

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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

func (api *Api) RunContainer(imageName string) (string, error) {
	ctx := context.Background()

	resp, err := api.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Tty:   true,
	}, nil, nil, nil, "")
	if err != nil {
		api.logger.Error(err.Error())
		return "", err
	}

	if err := api.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		api.logger.Error(err.Error())
		return "", err
	}

	return resp.ID, nil
}

func (api *Api) StopContainer(containerId string) error {
	ctx := context.Background()

	err := api.cli.ContainerStop(ctx, containerId, nil)
	if err != nil {
		api.logger.Error(err.Error())
		return err
	}

	return nil
}

func (api *Api) ContainerLogs(containerId string) (string, error) {
	ctx := context.Background()

	out, err := api.cli.ContainerLogs(ctx, containerId, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		api.logger.Error(err.Error())
		return "", err
	}

	buf := new(strings.Builder)
	_, err = stdcopy.StdCopy(buf, nil, out)
	if err != nil {
		_, err = io.Copy(buf, out)
		if err != nil {
			api.logger.Error(err.Error())
			return "", err
		}
	}

	return buf.String(), nil
}
