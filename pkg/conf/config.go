package conf

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

const (
	DmgMntPoint  = "/Volumes/evidence_image_dmg/"
	CmdXmount    = "/usr/local/bin/xmount"
	CmdDiskUtil  = "/usr/sbin/diskutil"
	CmdPlUtil    = "/usr/bin/plutil"
	CmdHdiUtil   = "/usr/bin/hdiutil"
	CmdMountApfs = "/sbin/mount_apfs"
)
var SupportFileType = []string{"ewf", "raw", "dmg"}

type Config struct {
	ImgPath     string
	VaultPW     string
	FileType    string
	FlagUnmount bool
}

func (config *Config) CheckFileType() (bool, error) {

	for _, ft := range SupportFileType {
		if ft == config.FileType {
			return true, nil
		}
	}
	return false, errors.New("[-] Not support such file type")
}

func (config *Config) CheckImgFile() (bool, error) {
	fileInfo, err := os.Stat(config.ImgPath)
	if err != nil {
		return false, err
	}

	if fileInfo.IsDir() {
		return false, errors.New("[-] Directory is set")
	}

	ext := filepath.Ext(config.ImgPath)[1:]
	e01Regs := regexp.MustCompile(`(?i)E\d{2}`)
	dmgRegs := regexp.MustCompile(`(?i)dmg`)
	rawRegs := regexp.MustCompile(`(?i)raw`)
	if !(e01Regs.MatchString(ext) || dmgRegs.MatchString(ext) || rawRegs.MatchString(ext)) {
		return false, errors.New("[-] File it not supported")
	}
 	return true, nil
}