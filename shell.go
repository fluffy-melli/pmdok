package main

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/fluffy-melli/pmdok/docker"
)

func handleArgs(client *client.Client, args []string) {
	if len(args) <= 1 || args[1] == "help" {
		printHelp()
		return
	}

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
			printHelp()
			return
		}
	}
}

func printHelp() {
	fmt.Println(`Usage:
  pmdok list
	- Retrieves the list from Docker.

  pmdok pull <image (e.g., ubuntu)>
	- Downloads a Docker image.

  pmdok new <image (e.g., ubuntu)> <name> <commands to run>
	- Creates a new Docker container. The root is the current directory.

  pmdok start <name>
	- Starts a Docker container.

  pmdok stop <name>
	- Stops a Docker container.

  pmdok del <name>
	- Deletes a Docker container.`)
}

func main() {
	client := docker.NewClient()
	handleArgs(client, os.Args)
}
