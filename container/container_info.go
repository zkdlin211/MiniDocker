package container

import (
	"MiniDocker/constant"
	"MiniDocker/log"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type ContainerInfo struct {
	// PID of the container's init process on its host
	Pid string `json:"pid"`
	// container id
	Id string `json:"id"`
	// container name
	Name string `json:"name"`
	// Command executed by the init process of the container
	Command string `json:"command"`
	// create time of the container
	CreateTime string `json:"createTime"`
	// container currect status
	Status string `json:"status"`
}

func NewContainerInfo(Pid int, commandArr []string, name string, id string) *ContainerInfo {
	return &ContainerInfo{
		Id:         id,
		Pid:        strconv.Itoa(Pid),
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     constant.RUNNING,
		Name:       name,
		Command:    strings.Join(commandArr, ""),
	}
}

func GetContainerInfo(file fs.DirEntry) (*ContainerInfo, error) {
	containerName := file.Name()
	configFileDir := fmt.Sprintf(constant.DefaultInfoLocation, containerName)
	configFileDir = filepath.Join(configFileDir, constant.ConfigFile)
	content, err := os.ReadFile(configFileDir)
	if err != nil {
		log.Errorf("[ContainerInfo.GetContainerInfo] error reading container "+
			"info file %s, %v", configFileDir, err)
		return nil, err
	}
	containerInfo := &ContainerInfo{}
	if err := json.Unmarshal(content, containerInfo); err != nil {
		log.Errorf("[ContainerInfo.GetContainerInfo] json unmarshal error: %v", err)
		return nil, err
	}
	return containerInfo, nil
}

// Record ContainerInfo instance to configuration file using Json serialization
func (this *ContainerInfo) Record() error {
	jsonBytes, err := json.Marshal(this)
	if err != nil {
		log.Errorf("[ContainerInfo.Record] error serialize ContainerInfo to json: %v", err)
		return err
	}
	recordUrl := fmt.Sprintf(constant.DefaultInfoLocation, this.Name)
	if err := os.MkdirAll(recordUrl, 0622); err != nil {
		log.Errorf("[ContainerInfo.Record] error mkdir %s, %v", recordUrl, err)
		return err
	}
	recordUrl = filepath.Join(recordUrl, constant.ConfigFile)
	file, err := os.Create(recordUrl)
	if err != nil {
		log.Errorf("[ContainerInfo.Record] error create file %s, %v", recordUrl, err)
		return err
	}
	if _, err = file.WriteString(string(jsonBytes)); err != nil {
		log.Errorf("[ContainerInfo.Record] error write config to file %s, %v", recordUrl, err)
		return err
	}
	return nil
}

func (this *ContainerInfo) DeleteInfo() {
	recordUrl := fmt.Sprintf(constant.DefaultInfoLocation, this.Name)
	if err := os.RemoveAll(recordUrl); err != nil {
		log.Errorf("[ContainerInfo.DeleteInfo] error remove container info file %s, %v", recordUrl, err)
	}
}

func PrintToStdout(infos []*ContainerInfo) {
	writer := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprintf(writer, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, info := range infos {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
			info.Id,
			info.Name,
			info.Pid,
			info.Status,
			info.Command,
			info.CreateTime,
		)
	}
	if err := writer.Flush(); err != nil {
		log.Errorf("[ContainerInfo.PrintToStdout] writer Flush error: %v", err)
	}
}
