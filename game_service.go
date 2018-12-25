/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-11-27 20:08:05
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-20 14:59:50
 * @Description: matchvs game server , the main module for start or stop server
 */

package matchvs

import (
	"commonlibs/errors"
	"commonlibs/servers"
	"os"
	"strings"

	conf "github.com/matchvs/gameServer-go/src/config"
	"github.com/matchvs/gameServer-go/src/defines"
	"github.com/matchvs/gameServer-go/src/log"
	"github.com/matchvs/gameServer-go/src/message"
)

var (
	GsConfig *conf.GSConfig
)

// 业务接口类型
type BaseGsHandler interface {
	message.IHandler
	SetPushHandler(PushHandler)
}

type PushHandler interface {
	// PushEvent 发送房间消息
	// req 要推送的消息
	PushEvent(req *defines.MsPushEventReq) error
	// JoinOver 关闭房间
	// gameID : 游戏ID
	// roomID ：房间ID
	JoinOver(gameID uint32, roomID uint64)
	// JoinOpen 打开房间
	// gameID : 游戏ID
	// roomID ：房间ID
	JoinOpen(gameID uint32, roomID uint64)
	// KickPlayer 踢除指定玩家
	// destID : 要踢除的玩家
	// roomID : 房间ID
	KickPlayer(destID uint32, roomID uint64)
	// GetRoomDetail 获取房间详细信息
	// gameID : 游戏ID
	// roomID ：房间ID
	// latestWatcherNum : 获取最新观战人数的房间
	GetRoomDetail(gameID, latestWatcherNum uint32, roomID uint64)
	// SetRoomProperty 设置房间属性
	// gameID : 游戏ID
	// roomID ：房间ID
	// roomProperty : 房间属性
	SetRoomProperty(gameID uint32, roomID uint64, roomProperty string)
	// CreateRoom 主动创建房间
	// crtm ： 创建房间的参数信息
	// 返回类型 MsCreateRoomRsp 是创建后的状态信息
	CreateRoom(crtm *defines.MsCreateRoomReq) (*defines.MsCreateRoomRsp, error)
	// TouchRoom 设置房间的存活时间
	// 返回 200 表示成功
	// gameID : 游戏ID
	// ttl 	  : 空房间存活时长(房间没有任何玩家的情况)，单位秒，最大取值 86400 秒（1天）
	// roomID : 房间ID
	TouchRoom(gameID, ttl uint32, roomID uint64) (uint32, error)
	// DestroyRoom 主动销毁房间，可以销毁任意房间，如果房间中有人，就会把人剔出房间
	// 返回 200 表示成功
	// gameID : 游戏ID
	// roomID ：房间ID
	DestroyRoom(gameID uint32, roomID uint64) (uint32, error)
	// SetFrameSyncRate 设置帧同步参数
	// gameID : 游戏ID
	// frameRate : 帧率（0到20，且能被1000整除）
	// enableGS GameServer是否参与帧同步（0：不参与；1：参与）
	// roomID : 要设置帧同步的房间ID
	SetFrameSyncRate(setinfo *defines.MsSetFrameSyncRateReq)
	// FrameBroadcast 发送帧同步数据给 游戏房间服务
	// gameID : 游戏ID
	// operation : 数据处理方式 0：只发客户端；1：只发GS；2：同时发送客户端和GS
	// roomID : 房间ID
	// cpProto : 要发送的数据
	FrameBroadcast(gameID uint32, operation int32, roomID uint64, cpProto []byte)

	// GetOffLineCacheData 获取断线重新连接后的缓存数据
	// gameID : 游戏ID
	// cacheMS : 想要获取的毫秒数(-1表示获取所有缓存数据，该字段的赋值上限为1小时)
	// roomID : 房间ID
	GetOffLineCacheData(gameID uint32, roomID uint64, cacheMS int32) error
}

// 初始化读取配置
func initConfig(confFile string) {
	var (
		err error
	)
	//获取命令行参数
	args := conf.NewTerminalCmd()
	if strings.TrimSpace(confFile) != "" {
		args.ConfFile = confFile
	}
	// 获取配置文件参数
	GsConfig, err = conf.NewGsConfig(args.ConfFile)
	//log.LogD("configuration file read success %s", args.ConfFile)
	if err != nil {
		panic(err)
	}
	//设置日志级别
	if GsConfig.Log != nil {
		log.LogI("log level to set [%s]", GsConfig.Log.Level)
		log.SetLevel(GsConfig.Log.Level)
	}
}

// game server main struct which use to start and stop server
type GameServer struct {
	handler BaseGsHandler
	adaptor *message.GSAdaptor
	roomMg  *servers.RoomManager
	push    *message.PushManager
	server  *servers.StreamServer
}

// NewGameServer is game server manager , you can start or stop to use it
// BaseGsHandler is base handler class , you neet to implement BaseGsHandler
// confFile the configuration file path (include the name). if confFile value is "" ,default read "./conf/config.toml"
func NewGameServer(hd BaseGsHandler, confFile string) (g *GameServer) {
	initConfig(confFile)
	g = new(GameServer)
	g.handler = hd
	g.adaptor = message.NewGSAdaptor(g.handler)
	g.server = servers.NewStreamServer(GsConfig.Server.Host, g.adaptor, 10)
	g.roomMg = servers.NewRoomManager(GsConfig.RoomManage)
	g.push = message.NewPushManager(g.adaptor, g.roomMg)
	return
}

// 启动 gameServer 服务
func (g *GameServer) Start() (err error) {
	defer errors.PrintPanicStack()
	register := servers.NewRegister(GsConfig.Register)
	register.Run()
	if err := g.server.Start(); err != nil {
		register.Stop()
		g.Stop()
		log.LogE("gameServer start err %v", err)
		os.Exit(-1)
	}
	register.Stop()
	return nil
}

// 停止 gameServer 服务
func (g *GameServer) Stop() {
	g.roomMg.Stop()
	g.server.Stop()
}

func (g *GameServer) GetPushHandler() PushHandler {
	if g.push == nil {
		g.push = message.NewPushManager(g.adaptor, g.roomMg)
	}
	return g.push
}
