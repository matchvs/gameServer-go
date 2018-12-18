/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-11-28 14:30:33
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-18 17:11:33
 * @Description: matchvs game server example
 */

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/matchvs/gameServer-go"
	"github.com/matchvs/gameServer-go/example/app"
)

//程序函数入口
func main() {
	// handler 实现接口 matchvs.BaseGsHandler
	handler := app.NewApp()
	// 创建 gameServer
	gsserver := matchvs.NewGameServer(handler, "")
	// 设置消息推送
	handler.SetPushHandler(gsserver.GetPushHandler())
	// 启动 gameSever 服务
	go gsserver.Start()
	//检测系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	processStr := <-sigCh
	gsserver.Stop()
	fmt.Printf("service close  %v", processStr)
}
