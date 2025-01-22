package main

import (
	"os"

	"github.com/docker/docker/client"
	"github.com/fluffy-melli/pmdok/docker"
)

func handleArgs(client *client.Client, args []string) {
	args = args[1:]
	for i, arg := range args {
		switch arg {
		case "del":
			docker.RemoveContainer(client, args[i+1])
		case "start":
			docker.StartContainer(client, args[i+1])
		case "stop":
			docker.StopContainer(client, args[i+1])
		case "pull":
			docker.PullImage(client, args[i+1])
		case "new":
			executable, err := os.Getwd()
			if err != nil {
				docker.Error(docker.ERROR, "Could not get executable path: %s", err)
			}
			docker.CreateAndStartContainer(client, args[i+1], args[i+2], executable, args[i+3:])
		case "list":
			docker.ContainerList(client)
		default:
			return
		}
	}
}

func main() {
	client := docker.NewClient()
	handleArgs(client, os.Args)
}
