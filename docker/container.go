package docker

import (
	"context"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
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
	if len(containers) == 0 {
		Error(ERROR, "\033[31mno containers found\033[0m\n")
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

func GetAllLogs(cli *client.Client) {
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		Error(ERROR, "\033[31mlist containers:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	if len(containers) == 0 {
		Error(ERROR, "\033[31mno containers found\033[0m\n")
	}
	for _, info := range containers {
		logs, err := cli.ContainerLogs(context.Background(), info.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			Error(ERROR, "\033[31mget logs for container:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
			return
		}
		defer logs.Close()
		buf := make([]byte, 1024)
		for {
			n, err := logs.Read(buf)
			if err != nil {
				break
			}
			os.Stdout.Write(buf[:n])
		}
	}
}

func GetLogs(cli *client.Client, containerID string) {
	logs, err := cli.ContainerLogs(context.Background(), containerID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		Error(ERROR, "\033[31mget logs for container:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	defer logs.Close()
	buf := make([]byte, 1024)
	for {
		n, err := logs.Read(buf)
		if err != nil {
			break
		}
		os.Stdout.Write(buf[:n])
	}
}

func PullImage(cli *client.Client, imagename string) {
	out, err := cli.ImagePull(context.Background(), imagename, image.PullOptions{})
	if err != nil {
		Error(ERROR, "\033[31mpull image:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	defer out.Close()
	buf := make([]byte, 1024)
	for {
		_, err := out.Read(buf)
		if err != nil {
			Error(ERROR, "\033[31mread image pull:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
			break
		}
		//os.Stdout.Write(buf[:n])
	}
	Print(INFO, "pulled image: \033[1m\033[37m\033[42m%s\033[0m\n", imagename)
}

func CreaftContainer(cli *client.Client, image string, name string, cmd []string) {
	config := &container.Config{
		Image: image,
		Cmd:   cmd,
	}
	resp, err := cli.ContainerCreate(context.Background(), config, nil, nil, nil, name)
	if err != nil {
		Error(ERROR, "\033[31mcreate container:\033[0m \033[1m\033[37m\033[41m%v\033[0m\n", err)
		return
	}
	Print(INFO, "created container: \033[1m\033[37m\033[42m%s\033[0m\n", resp.ID)
}
