package docker

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func NewClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		Error(ERROR, "\033[31mcreate docker client:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return nil
	}
	return cli
}

func ContainerList(cli *client.Client) {
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		Error(ERROR, "\033[31mlist containers:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	for _, container := range containers {
		status := container.State
		created := time.Unix(container.Created, 0)
		duration := time.Since(created)
		if status == "running" {
			Print(INFO, "\033[1m\033[37m\033[42m%s\033[0m (running for \033[1m\033[34m\033[43m%s\033[0m)\n", container.Names[0][1:], duration)
		} else {
			Print(ERROR, "\033[1m\033[37m\033[41m%s\033[0m (stopped \033[1m\033[34m\033[43m%s\033[0m ago)\n", container.Names[0][1:], duration)
		}
	}
}

func PullImage(cli *client.Client, imageName string) {
	start := time.Now().UnixMilli()
	reader, err := cli.ImagePull(context.Background(), imageName, image.PullOptions{})
	if err != nil {
		Error(ERROR, "\033[31mpull image:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	defer reader.Close()
	end := time.Now().UnixMilli()
	Print(INFO, "downloaded image: \033[1m\033[37m\033[42m%s / %fms\033[0m\n", imageName, float64(end-start)/1000)
}

func CreateAndStartContainer(cli *client.Client, imageName string, containerName string, workingDir string, cmd []string) {
	absWorkingDir, err := filepath.Abs(workingDir)
	if err != nil {
		Error(ERROR, "\033[31mresolve absolute working directory:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	if _, err := os.Stat(absWorkingDir); os.IsNotExist(err) {
		Error(ERROR, "\033[31mworking directory does not exist:\033[0m \033[1m\033[37m\033[41m%s\033[0m\n", absWorkingDir)
		return
	}
	containerConfig := &container.Config{
		WorkingDir: strings.ReplaceAll(absWorkingDir, "\\", "/"),
		Cmd:        cmd,
	}
	resp, err := cli.ContainerCreate(context.Background(), containerConfig, nil, nil, &v1.Platform{OS: "linux"}, containerName)
	if err != nil {
		Error(ERROR, "\033[31mcreate container:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	if err != nil {
		Error(ERROR, "\033[31mstart container:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	Print(INFO, "create container: \033[1m\033[37m\033[42m%s\033[0m\n", containerName)
}

func StartContainer(cli *client.Client, containerID string) {
	err := cli.ContainerStart(context.Background(), containerID, container.StartOptions{})
	if err != nil {
		Error(ERROR, "\033[31mstart container:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	Print(INFO, "started container: \033[1m\033[37m\033[42m%s\033[0m\n", containerID)
}

func StopContainer(cli *client.Client, containerID string) {
	timeout := 10
	err := cli.ContainerStop(context.Background(), containerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		Error(ERROR, "\033[31mstop container:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	Print(INFO, "stopped container: \033[1m\033[37m\033[42m%s\033[0m\n", containerID)
}

func RemoveContainer(cli *client.Client, containerID string) {
	err := cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true})
	if err != nil {
		Error(ERROR, "\033[31mremove container:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	Print(INFO, "removed container: \033[1m\033[37m\033[42m%s\033[0m\n", containerID)
}
