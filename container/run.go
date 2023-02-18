package container

import (
	"MiniDocker/cgroup"
	"MiniDocker/cgroup/subsystem"
	"MiniDocker/filesystem/aufs"
	"MiniDocker/log"
	"github.com/google/uuid"
	"os"
	"strings"
)

func StartContainer(tty bool, command []string, resConf *subsystem.ResourceConfig, volume string, containerName string) {
	log.Debugf("[container.StartContainer] starting container")
	id := uuid.New().String()
	if containerName == "" {
		containerName = id
	}
	parent, writePipe := NewParentProcess(tty, volume, containerName)
	if err := parent.Start(); err != nil {
		log.Error("[run.StartContainer] running parent process in container failed: %v", err.Error())
	}
	// create a corresponding containerInfo and record it to configuration file
	containerInfo := NewContainerInfo(parent.Process.Pid, command, containerName, id)
	if err := containerInfo.Record(); err != nil {
		return
	}
	// cgroup logic
	// use minidocker-cgroup as cgroup name
	// create cgroup manager
	manager := cgroup.NewCgroupManager("minidocker-cgroup")
	// configure the resource limit by parameters
	if err := manager.Set(resConf); err != nil {
		panic(err)
	}
	// Add the container process to the cgroup corresponding to each mounted subsystem
	if err := manager.Apply(parent.Process.Pid); err != nil {
		panic(err)
	}

	sendInitCommand(command, writePipe)

	if tty {
		if err := parent.Wait(); err != nil {
			log.Errorf("[run.StartContainer] container parent process waiting error: %v", err.Error())
			return
		}
		containerInfo.DeleteInfo()
		log.Info("[main] Remove container recourse")
		manager.Release()
		aufs.DeleteWorkSpace(volume, containerName)
	} else {
		// non tty, do nothing, init process(pid = 1) will adopt this process
	}
}

func sendInitCommand(command []string, writePipe *os.File) {
	cmd := strings.Join(command, " ")
	log.Debugf("[sendInitCommand] command: %v", cmd)
	if _, err := writePipe.WriteString(cmd); err != nil {
		panic(err)
	}
	writePipe.Close()
}
