package demo2

import (
	"MiniDocker/log"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

// The root directory location of the memory subsystem hierarchy that is mounted
const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		log.Debugf("current pid %d", syscall.Getpid())
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Printf("fork pid: %v\n", cmd.Process.Pid)

	os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "test-memory-limit"), 0755)
	os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "test-memory-limit", "tasks"),
		[]byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "test-memory-limit", "memory.limit_in_bytes"),
		[]byte("100m"), 0644)
	cmd.Process.Wait()
}
