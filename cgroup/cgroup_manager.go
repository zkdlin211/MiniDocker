package cgroup

import (
	"MiniDocker/cgroup/subsystem"
	"MiniDocker/log"
)

// CgroupManager will control various subsystems and link them with the containers
type CgroupManager struct {

	// hierarchy path for cgroup, i.e. path of the created cgroup path relative to each root cgroup path
	Path string

	Resource *subsystem.ResourceConfig
}

// Set the cgroup resource limit in each subsystem mount
func (this *CgroupManager) Set(res *subsystem.ResourceConfig) error {
	for _, sub := range subsystem.SubsystemIns {
		if err := sub.Set(this.Path, res); err != nil {
			log.Errorf("[CgroupManager.Set] error setting cgroup resource limit: %v", err.Error())
			return err
		}
	}
	return nil
}

func (this *CgroupManager) Apply(pid int) error {
	for _, sub := range subsystem.SubsystemIns {
		if err := sub.Apply(this.Path, pid); err != nil {
			log.Errorf("[CgroupManager.Apply] error writing pid to cgroup file: %v", err)
			return err
		}
	}
	return nil
}

// release all cgroups under each subsystem
func (this *CgroupManager) Release() error {
	for _, sub := range subsystem.SubsystemIns {
		if err := sub.Remove(this.Path); err != nil {
			return err
		}
	}
	return nil
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}
