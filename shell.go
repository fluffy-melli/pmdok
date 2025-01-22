package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/docker/docker/client"
	"github.com/fluffy-melli/pmdok/docker"
)

func handleArgs(client *client.Client, args []string) {
	parser := argparse.NewParser("pmdok", "Docker management tool")
	listCmd := parser.NewCommand("list", "Retrieves the list from Docker")
	pullCmd := parser.NewCommand("pull", "Downloads a Docker image")
	pullImage := pullCmd.String("i", "image", &argparse.Options{Required: true, Help: "Docker image to pull"})
	newCmd := parser.NewCommand("new", "Creates a new Docker container")
	newImage := newCmd.String("i", "image", &argparse.Options{Required: true, Help: "Docker image to use"})
	newName := newCmd.String("n", "name", &argparse.Options{Required: true, Help: "Name of the container"})
	newCommands := newCmd.String("c", "commands", &argparse.Options{Required: true, Help: "Commands to run in the container"})
	startCmd := parser.NewCommand("start", "Starts a Docker container")
	startName := startCmd.String("n", "name", &argparse.Options{Required: true, Help: "Name of the container"})
	stopCmd := parser.NewCommand("stop", "Stops a Docker container")
	stopName := stopCmd.String("n", "name", &argparse.Options{Required: true, Help: "Name of the container"})
	delCmd := parser.NewCommand("del", "Deletes a Docker container")
	delName := delCmd.String("n", "name", &argparse.Options{Required: true, Help: "Name of the container"})

	err := parser.Parse(args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	switch {
	case listCmd.Happened():
		docker.ContainerList(client)
	case pullCmd.Happened():
		docker.PullImage(client, *pullImage)
	case newCmd.Happened():
		executable, err := os.Getwd()
		if err != nil {
			docker.Error(docker.ERROR, "Could not get executable path: %s\n", err)
		}
		docker.CreateAndStartContainer(client, *newImage, *newName, executable, strings.Split((*newCommands), " "))
	case startCmd.Happened():
		docker.StartContainer(client, *startName)
	case stopCmd.Happened():
		docker.StopContainer(client, *stopName)
	case delCmd.Happened():
		docker.RemoveContainer(client, *delName)
	default:
		return
	}
}

func main() {
	client := docker.NewClient()
	handleArgs(client, os.Args)
}
