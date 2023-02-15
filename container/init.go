package container

import (
	"MiniDocker/log"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func StartContainerInitProcess(cmd string, args []string) error {
	log.Info("[container.StartContainerInitProcess] Starting container process")
	cmdArr := readUserCommand()
	if cmdArr == nil || len(cmdArr) == 0 {
		return fmt.Errorf("[container.StartContainerInitProcess] error starting container init process, " +
			"user command is nil")
	}

	path, err := exec.LookPath(cmdArr[0])
	if err != nil {
		log.Errorf("[container.StartContainerInitProcess] exec look path error: %v", err)
		return err
	}
	log.Infof("[container.StartContainerInitProcess] find path %s", path)
	if err := syscall.Exec(path, cmdArr[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Errorf("[container.readUserCommand] error reading user command: read pipe error: %v", err.Error())
		return nil
	}
	return strings.Split(string(msg), " ")
}
