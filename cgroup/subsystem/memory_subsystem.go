package subsystem

import (
	"MiniDocker/log"
	"os"
	"path"
	"strconv"
)

type MemorySubsystem struct {
}

func (this *MemorySubsystem) Name() string {
	return "memory"
}

// Set the memory resource limit of the cgroup corresponding to cgroupPath
func (this *MemorySubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	var subsysCgrouPath string
	var err error
	if subsysCgrouPath, err = GetCgroupPath(this.Name(), cgroupPath, true); err != nil {
		log.Errorf("[MemorySubsystem.Set] error when setting the memory resource limit: %v", err.Error())
		return err
	}
	if err = os.WriteFile(path.Join(subsysCgrouPath, "memory.limit_in_bytes"),
		[]byte(res.MemoryLimit), 0644); err != nil {
		log.Errorf("[MemorySubsystem.Set] error when writing the file: %v", err.Error())
		return err
	}
	return nil
}

// Add a process to the cgroup corresponding to cgroupPath
// `echo $(pid) >> /sys/fs/cgroup/memory/tasks`
func (this *MemorySubsystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(this.Name(), cgroupPath, false); err == nil {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			log.Errorf("[MemorySubsystem.Apply] error writing pid to cgroup file: %v", err)
			return err
		}
		return nil
	} else {
		return err
	}
}

// remove the corresponding cgroup
func (this *MemorySubsystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(this.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}
