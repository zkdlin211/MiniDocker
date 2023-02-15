package main

import (
	"MiniDocker/cgroup"
	"MiniDocker/cgroup/subsystem"
	"MiniDocker/container"
	"MiniDocker/log"
	"os"
	"strings"
)

func StartContainer(tty bool, command []string, resConf *subsystem.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if err := parent.Start(); err != nil {
		log.Error("[run.StartContainer] running parent process in container failed: %v", err.Error())
	}
	// cgroup logic
	// use minidocker-cgroup as cgroup name
	// create cgroup manager
	manager := cgroup.NewCgroupManager("minidocker-cgroup")
	defer manager.Release()
	// configure the resource limit by parameters
	if err := manager.Set(resConf); err != nil {
		panic(err)
	}
	// Add the container process to the cgroup corresponding to each mounted subsystem
	if err := manager.Apply(parent.Process.Pid); err != nil {
		panic(err)
	}

	sendInitCommand(command, writePipe)

	if err := parent.Wait(); err != nil {
		log.Error("[run.StartContainer] container parent process waiting error: %v", err.Error())
		return
	}
	os.Exit(-1)
}

func sendInitCommand(command []string, writePipe *os.File) {
	cmd := strings.Join(command, " ")
	log.Debugf("[sendInitCommand] command: %v", cmd)
	if _, err := writePipe.WriteString(cmd); err != nil {
		panic(err)
	}
	writePipe.Close()
}
