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
	default:
		return
	}
}

func main() {
	client := docker.NewClient()
	handleArgs(client, os.Args)
}
