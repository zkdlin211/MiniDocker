package container

import (
	"MiniDocker/log"
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
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
	}
	process.ExtraFiles = []*os.File{read}
	return process, write
}
