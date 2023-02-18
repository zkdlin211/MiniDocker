package command

import (
	"MiniDocker/constant"
	"MiniDocker/container"
	"MiniDocker/log"
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var ListCommand = cli.Command{
	Name:  "ps",
	Usage: "list all running containers",
	Action: func(context *cli.Context) error {
		ListContainers()
		return nil
	},
}

func ListContainers() {
	url := fmt.Sprintf(constant.DefaultInfoLocation, "")
	// remove `/`
	url = url[:len(url)-1]
	files, err := os.ReadDir(url)
	if err != nil {
		log.Errorf("[ListCommand.ListContainers] error read dir %s, %v", url, err)
		return
	}
	var infos []*container.ContainerInfo
	for _, file := range files {
		info, err := container.GetContainerInfo(file)
		if err != nil {
			log.Errorf("[ListCommand.ListContainers] error get container [%s] info, %v", file.Name(), err)
			continue
		}
		infos = append(infos, info)
	}
	container.PrintToStdout(infos)
}
