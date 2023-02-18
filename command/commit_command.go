package command

import (
	"MiniDocker/commit"
	"fmt"
	"github.com/urfave/cli"
)

var CommitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container and package it to an image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		imageName := context.Args().Get(0)
		return commit.CommitContainer(imageName)
	},
}
