package main

import (
	"github.com/matchvs/gameServer-go"
	"github.com/matchvs/gameServer-go/example/app"
)

//程序函数入口
func main() {
	// 定义业务处理对象这个业务类需要 继承接口
	handler := &app.App{GameID: uint32(123)}
	// 创建 gameServer
	gsserver := matchvs.NewGameServer(handler)
	handler.SetPushHandler(gsserver.GetPushHandler())
	// 启动 gameSever 服务
	gsserver.Start()
}
