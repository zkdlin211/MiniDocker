package subsystem

import (
	"MiniDocker/log"
	"os"
	"path"
	"strconv"
)

type CpusetSubSystem struct {
}

func (this *CpusetSubSystem) Name() string {
	return "cpuset"
}

func (this *CpusetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(this.Name(), cgroupPath, true); err == nil {
		if res.CpuSet != "" {
			if err := os.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"),
				[]byte(res.CpuSet), 0644); err != nil {
				log.Errorf("[CpusetSubSystem.Set] error setting `cpuset.cpus`, %v", err.Error())
				return err
			}
		}
		return nil
	} else {
		return err
	}
}

func (this *CpusetSubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(this.Name(), cgroupPath, false); err == nil {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			log.Errorf("[CpusetSubSystem.Apply] error writing pid to cgroup file: %v", err)
			return err
		}
		return nil
	} else {
		return err
	}
}

func (this *CpusetSubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(this.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}
