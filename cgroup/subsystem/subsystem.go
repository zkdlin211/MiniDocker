package subsystem

// all Subsystems
var (
	SubsystemIns = []Subsystem{
		&CpuSubsystem{},
		&MemorySubsystem{},
		&CpusetSubSystem{},
	}
)

// ResourceConfig is used to pass the resource restriction configuration
type ResourceConfig struct {
	MemoryLimit string
	// cpu time slice weight
	CpuShare string
	// cpu core limit
	CpuSet string
}

type Subsystem interface {
	// return the name of subsystem (eg. cpu memory)
	Name() string
	// Set the resource limit for a cgroup on this Subsystem
	Set(cgroupPath string, res *ResourceConfig) error
	// Add a process to a cgroup
	Apply(cgroupPath string, pid int) error
	// remove a cgroup
	Remove(cgroupPath string) error
}
