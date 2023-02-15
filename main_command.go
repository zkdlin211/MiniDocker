package main

import (
	"MiniDocker/cgroup/subsystem"
	"MiniDocker/container"
	"MiniDocker/log"
	"fmt"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroup limit miniDocker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpu share limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	// function that executes by 'miniDocker run' command
	Action: func(context *cli.Context) error {
		log.Debugf("[%s]", context.Args())
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArr []string
		for _, arg := range context.Args() {
			cmdArr = append(cmdArr, arg)
		}
		tty := context.Bool("it")
		resConf := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}
		StartContainer(tty, cmdArr, resConf)
		return nil
	},
}
var initCommand = cli.Command{
	Name: "init",
	Usage: "Initialize the container process that runs user's process in the container. " +
		"Only for INTERNAL use. DO NOT call it outside.",
	Action: func(context *cli.Context) error {
		log.Info("[main.initCommand] Initializing the container")
		cmd := context.Args().Get(0)
		log.Infof("[main.initCommand] Input command: %s", cmd)
		err := container.StartContainerInitProcess(cmd, nil)
		return err
	},
}
