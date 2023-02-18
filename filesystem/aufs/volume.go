package aufs

import (
	"MiniDocker/constant"
	"MiniDocker/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func VolumeURLExtract(volume string) ([]string, bool) {
	volumeURLs := strings.Split(volume, ":")
	return volumeURLs, len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != ""
}

func MountVolume(volumeURLs []string) {
	// create host file dir
	parentUrl := volumeURLs[0]
	if err := os.MkdirAll(parentUrl, 0777); err != nil {
		log.Infof("[aufs.MountVolume] mkdir parent dir %s error, %v", parentUrl, err)
	}
	// create container file dir
	containerUrl := volumeURLs[1]
	containerVolumeUrl := filepath.Join(constant.MntUrl, containerUrl)
	if err := os.MkdirAll(containerVolumeUrl, 0777); err != nil {
		log.Infof("[aufs.MountVolume] mkdir container volume dir %s error, %v", containerVolumeUrl, err)
	}
	// mount the host file directory to the container mount point
	dirs := "dirs=" + parentUrl
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("[aufs.MountVolume] Mount volume failed. %v", err)
	}
}
