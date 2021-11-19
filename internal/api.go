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

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			runtime.Events.Emit("containerUpdate", api.GetContainers())
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

func (api *Api) containerLogs(containerId string) (string, error) {
	ctx := context.Background()

	out, err := api.cli.ContainerLogs(ctx, containerId, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true})
	if err != nil {
		api.logger.Error(err.Error())
		return "", err
	}

	buf := new(strings.Builder)
	_, err = stdcopy.StdCopy(buf, buf, out)
	if err != nil {
		api.logger.Errorf("Could not use stdCopy to get the logs streams, will try to use io.Copy, error: %e", err.Error())

		_, err = io.Copy(buf, out)
		if err != nil {
			api.logger.Error(err.Error())
			return "", err
		}
	}

	return buf.String(), nil
}

func getLastLogTimeStamp(logs string) string {
	lastLine := ""
	lines := strings.Split(logs, "\n")
	if len(lines) > 0 {
		if lines[len(lines)-1] == "" && len(lines) > 1 {
			// Last line is null, so we need to return the second last line
			lastLine = lines[len(lines)-2]
		} else {
			lastLine = lines[len(lines)-1]
		}

	}

	if lastLine != "" && len(lastLine) >= 30 {
		return lastLine[:30]
	}

	return ""
}

func (api *Api) ListenForContainerLogs(containerId string) {
	ticker := time.NewTicker(time.Second)

	api.runtime.Events.On("container:log:stop", func(params ...interface{}) {
		ticker.Stop()
	})

	lastLogTimeStamp := ""
	lastSentLogTimeStamp := ""

	for range ticker.C {
		logs, err := api.containerLogs(containerId)
		if err != nil {
			api.logger.Error(err.Error())
			continue
		}

		timeStamp := getLastLogTimeStamp(logs)
		if timeStamp != "" {
			lastLogTimeStamp = timeStamp
		}

		if lastLogTimeStamp != lastSentLogTimeStamp && lastLogTimeStamp != "" {
			lastSentLogTimeStamp = lastLogTimeStamp
			api.runtime.Events.Emit("container:log:new", logs)
		}
	}
}
