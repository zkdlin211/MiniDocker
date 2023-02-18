package command

import (
	"MiniDocker/container"
	"MiniDocker/log"
	"github.com/urfave/cli"
)

var InitCommand = cli.Command{
	Name: "init",
	Usage: "Initialize the container process that runs user's process in the container. " +
		"Only for INTERNAL use. DO NOT call it outside.",
	Action: func(context *cli.Context) error {
		log.Info("[main.initCommand] Initializing the container")
		cmd := context.Args().Get(0)
		log.Infof("[main.initCommand] Input command: %s", cmd)
		err := container.StartContainerInitProcess()
		return err
	},
}
