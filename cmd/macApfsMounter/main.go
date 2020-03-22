package main

import (
	"github.com/Recruit-CSIRT/macApfsMounter/pkg/conf"
	"github.com/Recruit-CSIRT/macApfsMounter/pkg/gui"
)

func main() {
	var config conf.Config
	gui.Run(&config)
}