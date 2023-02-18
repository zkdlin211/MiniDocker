package commit

import (
	"MiniDocker/constant"
	"MiniDocker/log"
	"os/exec"
)

func CommitContainer(imageName string) error {
	imageTar := "/root/" + imageName + ".tar"
	log.Infof("[commit.CommitContainer] image url: %s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C",
		constant.MntUrl, ".").CombinedOutput(); err != nil {
		log.Errorf("[commit.CommitContainer] tar folder %s error %v", constant.MntUrl, err)
		return err
	}
	return nil
}
