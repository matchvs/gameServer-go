package config

import (
	"flag"
	"os"
)

type TerminalCmd struct {
	ConfFile   string
	ServerName string
	LogLevel   string
}

//获取命令行参数
func NewTerminalCmd() (cmd *TerminalCmd) {
	cmd = new(TerminalCmd)
	var help bool
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&cmd.ConfFile, "c", "./conf/config.toml", "the configuration file name")
	flag.StringVar(&cmd.ServerName, "s", "GameServer", "the name for service")
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	return
}
