package aufs

import (
	"MiniDocker/constant"
	"MiniDocker/log"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func NewWorkspace(volume string, containerName string) {
	CreateReadOnlyLayer()
	CreateWriteLayer(containerName)
	CreateMountPoint(containerName)
	if volume != "" {
		volumeURLs, valid := VolumeURLExtract(volume)
		// eg. /root/volume:/containerVolume
		if valid {
			MountVolume(volumeURLs)
			log.Infof("[aufs.NewWorkspace] volumes: %q", volumeURLs)
		} else {
			log.Errorf("[aufs.NewWorkspace] invalid volumes: %q", volumeURLs)
		}
	}
}

func CreateReadOnlyLayer() {
	busyboxURL := filepath.Join(constant.RootUrl, "/busybox/")
	busyboxTarURL := filepath.Join(constant.RootUrl, "/busybox.tar")
	exist, err := PathExists(busyboxURL)
	if err != nil {
		log.Infof("[aufs.CreateReadOnlyLayer] Fail to judge whether dir %s exists. %v", busyboxURL, err)
	}
	if exist == false {
		if err := os.MkdirAll(busyboxURL, 0777); err != nil {
			log.Errorf("[aufs.CreateReadOnlyLayer] Mkdir dir %s error. %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("[aufs.CreateReadOnlyLayer] Untar dir %s error %v", busyboxURL, err)
		}
	}
}

func CreateWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(constant.FWriteLayerUrl, containerName)
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		log.Errorf("[aufs.CreateWriteLayer] Mkdir dir %s error. %v", writeURL, err)
	}
}

// mount -t aufs -o dirs=$(/path/to/root/directory)= \
// ro:$(/path/to/root/directory)/writeLayer:$(/path/to/root/directory/)busybox none $(/path/to/mount/point)
func CreateMountPoint(containerName string) {
	mntUrl := filepath.Join(constant.MntUrl, containerName)
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		log.Errorf("[aufs.CreateMountPoint] Mkdir dir %s error. %v", mntUrl, err)
	}
	//tmpWriteLayer := fmt.Sprintf(constant.WriteLayerUrl, containerName)
	dirs := "dirs=" + filepath.Join(constant.RootUrl, "/writeLayer") + ":" + filepath.Join(constant.RootUrl, "/busybox")
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("[aufs.CreateMountPoint] %v", err)
	}
}

// Delete the AUFS filesystem while container exit
func DeleteWorkSpace(volume string, containerName string) {
	if volumeUrls, valid := VolumeURLExtract(volume); valid {
		DeleteMountPointWithVolume(containerName, volumeUrls)
	} else {
		DeleteMountPoint(containerName)
	}
	DeleteWriteLayer(containerName)
}

func DeleteMountPoint(containerName string) {
	mntURL := filepath.Join(constant.MntUrl, containerName)
	// delete `/root/mnt/${containerUrl}`
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("[aufs.DeleteMountPoint] %v", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("[aufs.DeleteMountPoint] Remove dir %s error %v", mntURL, err)
	}
}

func DeleteMountPointWithVolume(containerName string, volumeUrls []string) {
	// umount `/root/mnt/${containerUrl}`
	mntUrl := filepath.Join(constant.MntUrl, containerName)
	containerUrl := filepath.Join(mntUrl, volumeUrls[1])
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("[aufs.DeleteMountPointWithVolume] error umount volume. %v", err)
	}
	// umount `/root/mnt`
	cmd = exec.Command("umount", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("[aufs.DeleteMountPointWithVolume] error umount mount point. %v", err)
	}
	// remove mount point
	if err := os.RemoveAll(mntUrl); err != nil {
		log.Errorf("[aufs.DeleteMountPointWithVolume] error removing mount point dir %s. %v", mntUrl, err)
	}
}

func DeleteWriteLayer(containerName string) {
	writeURL := filepath.Join(constant.WriteLayerUrl, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("[aufs.DeleteWriteLayer] Remove dir %s error %v", writeURL, err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
