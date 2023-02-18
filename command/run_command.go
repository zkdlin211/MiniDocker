package command

import (
	"MiniDocker/cgroup/subsystem"
	"MiniDocker/container"
	"MiniDocker/log"
	"fmt"
	"github.com/urfave/cli"
)

var RunCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroup limit miniDocker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
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
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
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
		detach := context.Bool("d")
		volume := context.String("v")
		if tty && detach {
			return fmt.Errorf("`-it` and `-d` parameter can not be applied at the same time")
		}
		resConf := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}
		log.Debugf("[command.RunCommand] start container")
		container.StartContainer(tty, cmdArr, resConf, volume, context.String("name"))
		return nil
	},
}
