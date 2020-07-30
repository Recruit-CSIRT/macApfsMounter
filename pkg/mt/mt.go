package mt

import (
	"errors"
	"fmt"
	"github.com/Recruit-CSIRT/macApfsMounter/pkg/conf"
	"github.com/Recruit-CSIRT/macApfsMounter/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func mount(mountPoint string, cmd []string) error {

	// makr mount point dir
	err := os.MkdirAll(mountPoint, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to make dir")
		return err
	}

	c := exec.Command(cmd[0], cmd[1:]...)
	if out, err := c.Output(); err != nil {
		return errors.New("Message: " + string(out) + ", Error: " + err.Error())
	}

	return nil
}

func stringInArray(s string, li []string) bool {
	for _, cs := range li {
		if s == cs {
			return true
		}
	}
	return false
}

func Run (config *conf.Config) error {
	var err error
	// stap 1 mount e01
	// E01 to E*
	var imgPaths []string
	filename := filepath.Base(config.ImgPath[:len(config.ImgPath) - len(filepath.Ext(config.ImgPath))])
	if config.FileType == "ewf" && strings.HasPrefix(filepath.Ext(config.ImgPath), ".E") && len(filepath.Ext(config.ImgPath)) == 4 {
		var pattern = filepath.Dir(config.ImgPath) + "/" + filename + ".E[0-9A-Z][0-9A-Z]"
		imgPaths, _ = filepath.Glob(pattern)
	}

	// xmount
	fmt.Println("[+] Make a mount dir to put dmg file.")
	// sudo xmount --in ewf david_lightman_system.E01 --out dmg /tmp/evidence_image_dmg/
	if config.FileType == "ewf" || config.FileType == "raw" {
		cmd := []string{ conf.CmdXmount, "--in", config.FileType}
		cmd = append(cmd, imgPaths...)
		cmd = append(cmd, "--out", "dmg", conf.DmgMntPoint)
		if err = mount(conf.DmgMntPoint, cmd); err != nil{
			return err
		}
	}

	// step 2 現在のディスクのリストを取得
	fmt.Println("[+] Check APFS Volume Disk before attaching the evidence.")
	diskListFirst := utils.NewDiskList()
	if err = diskListFirst.Set(); err != nil {
		return err
	}

	initDiskNames := []string {}
	for _, apfs := range diskListFirst.Containers {
		initDiskNames = append(initDiskNames, apfs.ContainerReference)
	}

	// step 3 dmgをマウントする
	// sudo hdiutil attach -nomount /Volumes/evidence_image_dmg/david_lightman_system.dmg
	fmt.Println("[+] Attach disk image.")

	// for dmg
	dmgPath := config.ImgPath

	// for ewf or raw
	if config.FileType == "ewf" || config.FileType == "raw" {
		dmgPath = filepath.Join(conf.DmgMntPoint, filename + ".dmg")
	}

	cmd := exec.Command(conf.CmdHdiUtil, "attach", "-nomount", dmgPath)
	if out, err := cmd.Output(); err != nil {
		return errors.New("Message: " + string(out) + ", Error: " + err.Error() )
	}

	// step 4 差分のディスクを取得
	fmt.Println("[+] Check APFS Volume Disk after the evidence is attached. ")
	diskListSecond := utils.NewDiskList()
	if err = diskListSecond.Set(); err != nil {
		return err
	}

	// step 5  & step 6
	attachedVols := []string{}
	for _, apfs := range diskListSecond.Containers {
		if stringInArray(apfs.ContainerReference, initDiskNames) {
			continue
		}

		for _, vol := range apfs.Volumes {

			if vol.Name == "VM" || vol.Name == "Preboot" || vol.Name == "Recovery" {
				continue
			}

			// step 5 filevault unlock
			devPath := filepath.Join("/dev", vol.DeviceIdentifier)
			volPath := filepath.Join("/Volumes", vol.DeviceIdentifier)

			if len(config.VaultPW) > 0 && vol.FileVault == true && vol.Encryption == true  {
				fmt.Println("[+] Unlock FileVault of " + devPath)
				c := exec.Command(conf.CmdDiskUtil, "apfs", "unlockVolume", devPath, "-passphrase", config.VaultPW, "-nomount")
				out, err := c.Output()
				if err != nil {
					fmt.Println("[-] Unlock looks like failed. Message: " + string(out) + ", Error: " + err.Error())
				} else {
					fmt.Println("[+] Unlocked " + devPath)
				}
			}

			// step 6 mount
			fmt.Println("[+] Mounting device: " + devPath)
			cmd := []string{conf.CmdMountApfs, "-o", "rdonly,noexec,noowners", devPath, volPath}
			if err = mount(volPath, cmd); err != nil {
				fmt.Println("[-] Failed to mount " + devPath + ". Error: " + err.Error())
				fmt.Println("[-] Try manually: " + strings.Join(cmd, " "))
			} else {
				fmt.Println("[+] Mounted: " + devPath + " -> " + volPath)
				attachedVols = append(attachedVols, volPath)
			}
		}
	}

	if len(attachedVols) == 0{
		return errors.New("[-] no mounted volumes")
	}

	return nil
}


func Unmount() error {
	/*
	unount
	$ diskutil unmount /Volumes/mnt/
	Volume Macintosh HD on disk7s1 unmounted

	$ diskutil eject /dev/disk7
	Disk /dev/disk7 ejected

	$ sudo diskutil unmount /Volumes/tmp
	or
	$ sudo hdiutil unmount -force tmp
	*/
	c := exec.Command(conf.CmdDiskUtil, "unmount", "force", conf.DmgMntPoint)
	out, err := c.Output();
	if err != nil {
		fmt.Println(string(out))
		return err
	}
	return nil
}
