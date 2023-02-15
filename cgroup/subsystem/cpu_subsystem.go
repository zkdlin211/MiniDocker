package subsystem

import (
	"MiniDocker/log"
	"os"
	"path"
	"strconv"
)

// CpuSubsystem is a struct for controlling cpu time slice weight
type CpuSubsystem struct {
}

func (this *CpuSubsystem) Name() string {
	return "cpu"
}

func (this *CpuSubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(this.Name(), cgroupPath, true); err != nil {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
			log.Errorf("[CpuSubsystem.Set] error setting `cpu.shares`, %v", err.Error())
			return err
		}
		return nil
	} else {
		return err
	}
}

func (this *CpuSubsystem) Apply(cgroupPath string, pid int) error {
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

func (this *CpuSubsystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(this.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}
