package main

import (
	"flag"
	"fmt"
	"github.com/Recruit-CSIRT/macApfsMounter/pkg/conf"
	"github.com/Recruit-CSIRT/macApfsMounter/pkg/mt"
	"os"
)

func init() {

	flag.StringVar(&config.ImgPath, "i", "", "set the img path.")
	flag.StringVar(&config.VaultPW, "p", "", "set the password of FileVault2 on evidence.")
	flag.BoolVar(&config.FlagUnmount, "u", false, "unmount option. only ewf and raw")

	flag.StringVar(&config.FileType, "t", "ewf", "select the file type. ewf(e01 file) or dmg.")
	//flag.StringVar(&config.FileType, "t", "ewf", "select the file type. ewf(e01 file), raw or dmg.")
}

var config conf.Config

func main() {

	flag.Parse()

	if config.FlagUnmount {
		if err := mt.Unmount(); err != nil {
			fmt.Println("[-] Failed to unmount. ", err.Error())
		} else{
			fmt.Println("[+] Success to unmount. ")
		}
		os.Exit(0)
	}

	if len(config.ImgPath) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Println("[+] Tool start.")

	if f, err := config.CheckFileType(); !f {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if f, err := config.CheckImgFile(); !f {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := mt.Run(&config); err != nil {
		fmt.Println("[-] Failed to mount. ", err.Error())
	} else {
		fmt.Println("[+] Success to mount. ")
	}
	fmt.Println("[+] Tool finish.")
}