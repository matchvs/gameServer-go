package test

import (
	"fmt"
	"testing"
	"time"

	matchvs "github.com/matchvs/gameServer-go"
	"github.com/matchvs/gameServer-go/example/app"
)

func Test_MainServer(t *testing.T) {
	// 定义业务处理对象这个业务类需要 继承接口
	handler := &app.App{}
	// 创建 gameServer
	gsserver := matchvs.NewGameServer(handler, "../../example/conf/config.toml.sample")
	handler.SetPushHandler(gsserver.GetPushHandler())

	go func() {
		// 启动 gameSever 服务
		gsserver.Start()
	}()

	for {
		select {
		case <-time.After(time.Second * 10):
			fmt.Println("结束服务")
			gsserver.Stop()
			time.Sleep(time.Second * 2)
			return
		}
	}
}
