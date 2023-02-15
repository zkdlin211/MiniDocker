package main

import (
	"MiniDocker/log"
	"github.com/urfave/cli"
	"os"
)

const usage = `MiniDoker is a simplified version of Docker. 
The purpose of this project if understand how Docker works. JUST FOR FUN.`

func main() {
	app := cli.NewApp()
	app.Name = "MiniDocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("[main] error when starting MiniDocker: %v", err.Error())
	}
}
