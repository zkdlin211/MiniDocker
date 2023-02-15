package subsystem

import (
	"MiniDocker/log"
	"bufio"
	"os"
	"path"
	"strings"
)

// example record (centos 7):
// 36 30 0:32 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime shared:13 - cgroup cgroup rw,memory
// The last option of above exampple is `rw,memory`.
// From this piece of information we know that the subsystem mounted is memory.
// Then the folder created in /sys/fs/cgroup/memory corresponding to the cgroup created can be used to limit the memory.
func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		log.Errorf("[utils.FindCgroupMountPoint] error opening `/proc/self/mountinfo`: %v", err)
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("[utils.FindCgroupMountPoint] error scanning `/proc/self/mountinfo`: %v", err)
		return ""
	}
	return ""
}

// Get the absolute path of cgroup in the file system
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountPoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil ||
		(autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
				log.Errorf("[utils.GetCgroupPath] error getting the absolute "+
					"path of cgroup in the file system: %v", err.Error())
				return "", err
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", err
	}
}
