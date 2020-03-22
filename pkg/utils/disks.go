package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/Recruit-CSIRT/macApfsMounter/pkg/conf"
)

type DiskList struct {
	Containers []Containers `json:"containers"`
}

type Containers struct {
	Fusion                  bool             `json:"Fusion",bool`
	ContainerReference      string           `json:"ContainerReference",string`
	APFSContainerUUID       string           `json:"APFSContainerUUID",string`
	PhysicalStores          []PhysicalStores `json:"PhysicalStores"`
	CapacityCeiling         int64            `json:"CapacityCeiling",int64`
	Volumes                 []Volumes        `json:"Volumes"`
	CapacityFree            int64            `json:"CapacityFree",int64`
	DesignatedPhysicalStore string           `json:"DesignatedPhysicalStore",string`
}

type PhysicalStores struct {
	DiskUUID         string `json:"DiskUUID",string`
	DeviceIdentifier string `json:"DeviceIdentifier",string`
	Size             int64  `json:"Size",int64`
}

type Volumes struct {
	Locked            bool     `json:"Locked",bool`
	APFSVolumeUUID    string   `json:"APFSVolumeUUID",string`
	CapacityQuota     int      `json:"CapacityQuota",int`
	DeviceIdentifier  string   `json:"DeviceIdentifier",string`
	CapacityReserve   int      `json:"CapacityReserve",int`
	CryptoMigrationOn bool     `json:"CryptoMigrationOn",bool`
	Name              string   `json:"Name",string`
	Encryption        bool     `json:"Encryption",bool`
	CapacityInUse     int64    `json:"CapacityInUse",int64`
	FileVault         bool     `json:"FileVault",string`
	Roles             []string `json:"Roles",string`
}

func NewDiskList() DiskList {
	return DiskList{}
}

func (dl *DiskList) Set() error {
	var buff bytes.Buffer
	var err error
	if err = getDiskList(&buff); err != nil {
		return err
	}

	if err := json.Unmarshal(buff.Bytes(), dl); err != nil {
		fmt.Println("[-] JSON Unmarshal error:", err)
		return err
	}
	return nil
}

func getDiskList(output *bytes.Buffer) error {
	cmds := []*exec.Cmd{
		exec.Command(conf.CmdDiskUtil, "apfs", "list", "-plist"),
		exec.Command(conf.CmdPlUtil, "-convert", "json", "-o", "-", "-"),
	}

	// stdout
	cmds[len(cmds)-1].Stdout = output

	// pipe
	var err error
	for i, cmd := range cmds {
		cmd.Stderr = os.Stderr
		if i > 0 {
			if cmds[i].Stdin, err = cmds[i-1].StdoutPipe(); err != nil {
				return err
			}
		}
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}
	return nil
}
