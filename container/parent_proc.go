package container

import (
	"MiniDocker/constant"
	"MiniDocker/filesystem/aufs"
	"MiniDocker/log"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func NewParentProcess(tty bool, volume string, containerName string) (*exec.Cmd, *os.File) {
	read, write, err := os.Pipe()
	if err != nil {
		log.Errorf("[container.NewParentProcess] New pipe error: %v", err.Error())
		return nil, nil
	}
	// starting the process by calling itself
	process := exec.Command("/proc/self/exe", "init")
	process.SysProcAttr = &syscall.SysProcAttr{
		// fork out a new process and use namespace to isolate the newly created process
		// from the external environment
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	// If the user specifies the -it parameter, the input/output of the current process needs to
	// be imported to the standard input/output
	if tty {
		process.Stdin = os.Stdin
		process.Stdout = os.Stdout
		process.Stderr = os.Stderr
	} else {
		infoUrl := fmt.Sprintf(constant.DefaultInfoLocation, containerName)
		if err := os.MkdirAll(infoUrl, 0622); err != nil {
			log.Errorf("[container.NewParentProcess] error mkdir %s, %v", infoUrl, err)
			return nil, nil
		}
		stdLogFilePath := filepath.Join(infoUrl, constant.ContainerLogFile)
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Errorf("[container.NewParentProcess] error create log file for %s, %v", containerName, err)
			return nil, nil
		}
		process.Stdout = stdLogFile
	}
	process.ExtraFiles = []*os.File{read}
	aufs.NewWorkspace(volume, containerName)
	process.Dir = fmt.Sprintf(constant.FMntUrl, containerName)
	return process, write
}
