package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

/**
 * GameServer 配置文件类型
 */
type GSConfig struct {
	Server     *ServerConf   `json:"server"`
	Log        *LogConf      `json:"log"`
	Register   *RegisterConf `json:"register"`
	RoomManage *RoomConf     `json:"room_manage"`
}

type ServerConf struct {
	Host string `json:"host"`
}

/**
 * 日志配置类型
 */
type LogConf struct {
	Level string `json:"level"`
}

type RegisterConf struct {
	Enable     bool   `json:"enable"`
	GameID     uint32 `json:"game_id"`
	SvcName    string `json:"svc_name"`
	PodName    string `json:"pod_name"`
	RemoteHost string `json:"remote_host"`
	RemotePort int32  `json:"remote_port"`
	LocalHost  string `json:"local_host"`
	LocalPort  uint32 `json:"local_port"`
}

type RoomConf struct {
	Enable     bool   `json:"enable"`
	SvcName    string `json:"svc_name"`
	PodName    string `json:"pod_name"`
	RemoteHost string `json:"remote_host"`
	RemotePort int32  `json:"remote_port"`
}

//读取配置文件并创建 GSConfig 对象
func NewGsConfig(fileName string) (gs *GSConfig, err error) {
	gs = new(GSConfig)
	if err = readConfig(fileName, gs); err != nil {
		fmt.Println("read configuration error: ", err)
		panic(err)
	}
	if gs.Server == nil {
		gs.Server = &ServerConf{
			Host: "0.0.0.0:12345",
		}
	}
	return
}

//读取 toml 文件
func readConfig(fileName string, v interface{}) (err error) {
	var (
		file *os.File
		data []byte
	)
	if file, err = os.Open(fileName); err != nil {
		return
	}

	if data, err = ioutil.ReadAll(file); err != nil {
		return
	}

	if err = toml.Unmarshal(data, v); err != nil {
		return
	}
	return
}
