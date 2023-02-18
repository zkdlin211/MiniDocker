package container

import (
	"MiniDocker/log"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func StartContainerInitProcess() error {
	log.Info("[container.StartContainerInitProcess] Starting container process")
	cmdArr := readUserCommand()
	if cmdArr == nil || len(cmdArr) == 0 {
		return fmt.Errorf("[container.StartContainerInitProcess] error starting container init process, " +
			"user command is nil")
	}
	setUpMount()
	// look up the absolute path of the command in the system PATH
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
	// uintptr(3) indicates the file descriptor with index of 3, which is the end of the pipe passed in
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Errorf("[container.readUserCommand] error reading user command: read pipe error: %v", err.Error())
		return nil
	}
	return strings.Split(string(msg), " ")
}

// setUpMount function sets up the mount points for the container process.
func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("[container.setUpMount] error getting pwd: %v", err)
		return
	}
	log.Debugf("[container.setUpMount] Current location is %s", pwd)
	if err = pivotRoot(pwd); err != nil {
		log.Error(err)
	}
	// MS_NOEXEC: This flag ensures that executable files cannot be executed from the mounted file system.
	// This is a security measure that prevents potentially harmful code from being executed from a mounted file system.
	// MS_NOSUID: This flag ensures that the setuid bit is not honored for executables in the mounted file system.
	// This is a security measure that prevents privilege escalation attacks.
	// MS_NODEV: This flag ensures that device files are not created in the mounted file system.
	// This is a security measure that prevents unauthorized access to hardware resources.
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// The /proc file system is a virtual file system that provides access to process and system information.
	// It is mounted by the mount command at boot time and is required for many system-level operations,
	// such as accessing process information, system statistics, and kernel configuration.
	if err = syscall.Mount("proc", "/proc",
		"proc", uintptr(defaultMountFlags), ""); err != nil {
		log.Errorf("[container.setUpMount] error mounting /proc: %v", err)
	}
	// The /dev file system is a virtual file system that provides access to device files. Device files represent
	// hardware devices and are used by applications to communicate with them. The /dev file system is typically
	// mounted at boot time and is required for many system-level operations that involve hardware access.
	// tmpfs is a temporary file system that stores its files in memory. The reason for using tmpfs is that
	// it provides a fast and efficient way to store temporary files, such as those created by hardware drivers.
	// By mounting /dev as a tmpfs file system, we can provide a fast and efficient way for hardware drivers to store
	// temporary files, without the need for an actual physical device
	if err = syscall.Mount("tmpfs", "/dev", "tmpfs",
		syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755"); err != nil {
		log.Errorf("[container.setUpMount] error mounting /dev: %v", err)
	}
}

/*
implementation of system call pivot_root
demo:

	mkdir /myrootfs/.pivot_root

	# Create a bind mount of the root file system directory to itself
	mount --bind /myrootfs /myrootfs

	# Change the root file system to the new directory and move the old root to a temporary directory
	pivot_root /myrootfs /myrootfs/.pivot_root

	# Change the current working directory to the new root directory
	cd /

	# Unmount the temporary directory
	umount /.pivot_root

	# Remove the temporary directory
	rmdir /.pivot_root
*/
func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("[container.pivotRoot] Mount rootfs o itself error: %v", err)
	}
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("[container.pivotRoot] create %s dir error: %v", pivotDir, err)
	}
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("[container.pivotRoot] syscall pivot_root error: %v", err)
	}
	// change current dir
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("[container.pivotRoot] syscall chdir / error: %v", err)
	}
	pivotDir = filepath.Join("/", ".pivot_root")
	// unmount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("[container.pivotRoot] syscall unmount %s error: %v", pivotDir, err)
	}
	// remove temporary dir
	return os.Remove(pivotDir)
}
