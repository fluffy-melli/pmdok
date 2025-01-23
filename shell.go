package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/docker/docker/client"
	"github.com/fluffy-melli/pmdok/docker"
)

func handleArgs(client *client.Client, args []string) {
	parser := argparse.NewParser("pmdok", "Docker management tool")
	listCmd := parser.NewCommand("list", "Retrieves the list from Docker")
	startCmd := parser.NewCommand("start", "Starts a Docker container")
	startName := startCmd.String("n", "name", &argparse.Options{Required: true, Help: "Name of the container"})
	stopCmd := parser.NewCommand("stop", "Stops a Docker container")
	stopName := stopCmd.String("n", "name", &argparse.Options{Required: true, Help: "Name of the container"})
	delCmd := parser.NewCommand("del", "Deletes a Docker container")
	delName := delCmd.String("n", "name", &argparse.Options{Required: true, Help: "Name of the container"})
	logsCmd := parser.NewCommand("log", "Retrieves the logs from a Docker container")
	logsCmdName := logsCmd.String("n", "name", &argparse.Options{Required: false, Help: "Name of the container"})
	creaftCmd := parser.NewCommand("create", "Creates a Docker container")
	creaftCmdName := creaftCmd.String("n", "name", &argparse.Options{Required: true, Help: "Name of the container"})
	creaftCmdImage := creaftCmd.String("i", "image", &argparse.Options{Required: true, Help: "Image of the container"})
	creaftCmdCmd := creaftCmd.StringList("c", "cmd", &argparse.Options{Required: true, Help: "Command of the container"})
	pullCmd := parser.NewCommand("pull", "Pulls a Docker image")
	pullCmdImage := pullCmd.String("i", "image", &argparse.Options{Required: true, Help: "Image to pull"})

	err := parser.Parse(args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	switch {
	case listCmd.Happened():
		docker.ContainerList(client)
	case startCmd.Happened():
		docker.StartContainer(client, *startName)
	case stopCmd.Happened():
		docker.StopContainer(client, *stopName)
	case delCmd.Happened():
		docker.RemoveContainer(client, *delName)
	case logsCmd.Happened():
		if *logsCmdName == "" {
			docker.GetAllLogs(client)
			return
		}
		docker.GetLogs(client, *logsCmdName)
	case pullCmd.Happened():
		docker.PullImage(client, *pullCmdImage)
	case creaftCmd.Happened():
		docker.CreaftContainer(client, *creaftCmdImage, *creaftCmdName, *creaftCmdCmd)
	default:
		return
	}
}

func main() {
	client := docker.NewClient()
	handleArgs(client, os.Args)
}
