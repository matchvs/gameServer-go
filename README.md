# Matchvs Golang 版 gameServer

使用 golang 版本的 gameServer 须知：

- 需要在本地安装好 golang 运行环境。如果还没有配置好 golang 环境的请先配置好。以下内容教程是默认你已经配置好并且熟悉 golang 开发环境。
- 在 Matchvs 官网创建了 gameServer。

[github 地址](https://github.com/matchvs/gameServer-go) 

## 新建 gameServer

前往 [matchvs 官网](http://www.matchvs.com/manage/gameServer) 新建一个 gameServer 服务，新建方法可以参考 [nodejs 版gameServer](https://doc.matchvs.com/QuickStart/GameServer-JavaScript) 版本的说明。



## 配置 gameServer 工程

#### 1、获取 matchvs golang 版本 gameServer 服务代码

使用如下命令获取拉取代码：

```shell
go get -u github.com/matchvs/gameServer-go
```

拉取成功后可以在你的 $GOPATH/src/github.com 目录下面看到 matchvs/gameServer-go 的服务代码。

在你的 `$GOPATH/src/github.com/ matchvs/gameServer-go` 目录下面有一个 example 文件，文件中的内容结构如下：

```
┌ app/  		# 业务代码
├ conf/ 		# 配置文件
├ main.go 		# 程序入口
├ Dockerfile 	# docker部署，这个文件不要修改
├ Makefile 		# make 文件内容不能修改
└ gsmeta 		# 数据源文件内容不能修改
```

#### 2、拉取 matchvs 官网的 gameServer 仓库

- 在你的 GOAPTH 目录下面使用 git 工具 拉取你在  [matchvs 官网](http://www.matchvs.com/manage/gameServer) 新建的 gameServer 仓库到本地。第一次拉取的仓库内容为空。如下示例：

```shell
  git clone https://git.matchvs.com/xxxxxxxxxx.git gogs-demo
```

- 复制 `$GOPATH/src/github.com/matchvs/gameServer-go/example/` 文件里面的内容到你刚刚拉下的仓库目录。

#### 3、修改配置文件

- 打开 conf/config.toml.sample 文件可以看到如下内容：我们刚开始只需要关心 [Server] 下面的 Host内容。

```toml
[Server]
Host = "0.0.0.0:36520"

[Log]
# Level's value can INFO、DEBUG、ERROR、WARNING, the default value is INFO
Level = "DEBUG"

[Register]
Enable = false
GameID = 123456
SvcName = "svc-abc"
PodName = "pod-abc"
RemoteHost = "192.168.8.1"
RemotePort = 9981
LocalHost = "192.168.8.2"
LocalPort = 30054

[RoomManage]
Enable = false
SvcName = "svc-abc"
PodName = "pod-abc"
RemoteHost = "192.168.8.1"
RemotePort = 9982
```

- 修改 Host 字段的端口号为你在 matchvs 官网创建 gameServer 后生成的端口号，例如你的 gameServer 端口号为 36000，那么你就修改 Host 字段为 0.0.0.0:36000

```toml
[Server]
Host = "0.0.0.0:36000"
......
```

- 重命名 conf/config.toml.sample 文件为 conf/config.toml，因为 gameServer 服务默认读取的配置文件是 conf/config.toml 。如果你没有修改配置文件名称在运行程序的时候会出现如下错误：

```shell
read configuration error:  open ./conf/config.toml: The system cannot find the file specified.
panic: open ./conf/config.toml: The system cannot find the file specified.

goroutine 1 [running]:
github.com/matchvs/gameServer-go/src/config.NewGsConfig(0x94f398, 0x12, 0x0, 0x0, 0x4004206bea0)
E:/Work/GOPATH/src/github.com/matchvs/gameServer-go/src/config/config.go:56 +0x1d2
........
exit status 2
```

- 修改和 main.go 文件 import 中的内容。在  main.go import 文件中我们看到如下内容:

```go
"github.com/matchvs/gameServer-go/example/app"
```

  这条引用需要改为你自己项目的引用路径，比如你现在的项目名称为 gogs-demo ，那么你需要修改为：

```go
"gogs-demo/app"
```

  > 注意：以上内容是 golang 包导入的知识，如果不太理解请阅读 golang 相关的知识。
  >
  >   如果没有修改 "github.com/src/matchvs/gameServer-go/example/app"  内容为你自己的项目内容，则项目运行将运行 example 中的 app.go

- 修改好 import 内容后就可以运行 第一个 gameServer 程序啦。运行内容大概如下：

```shell
  2018/11/30 16:26:43.63 game_service.go:55 ▶ INFO 001 log level to set [DEBUG]
  2018/11/30 16:26:43.63 room_manager.go:51 ▶ INFO 002 room manager is disable
  2018/11/30 16:26:43.63 dir_register.go:46 ▶ INFO 003 register is disable
  2018/11/30 16:26:43.63 stream_server.go:66 ▶ INFO 004 game server start success [0.0.0.0:36520]
```



## 本地调试

本地调试支持 windows、linux、mac 等多平台。这里就以使用较多的 windows 为例简单说明。

> **注意：** 本地调试模式，相当于 客户端是 alpha 环境，所以在使用SDK的时候，login 接口参数需使用 alpha

#### 1、下载配置 matchvs 命令行工具

[命令行工具](http://www.matchvs.com/serviceDownload) 下载 win 版本，提出 matchvs.exe 文件，并配置好环境变量。在终端输入 matchvs 看到如下结果说明配置成功。

```shell
Usages: matchvs COMMAND

The commands are:
        login           账号登录
        list            查询gameServer列表
        info            查询单个gameServer信息
        create          创建gameServer
        modify          修改gameServer
        delete          删除gameServer
        publish         发布gameServer到镜像仓库
        run             控制gameServer启停
        debug           开启本地调试，调试完成后在终端输入"quit"退出本地调试

Use "matchvs help [command]" for more information about a command.

```

#### 2、打开本地调试

使用 matchvs login 登录获取权限验证，登录步骤和信息如下：

```shell
matchvs login
Email or Mobile phone: 这里是你注册matchvs的邮箱或者手机号码
Password: 这个密码是你新建 gameServer后在官网生成的密码，不是你的账号密码
1 -- Matchvs
2 -- Cocos
3 -- Egret（白鹭）
4 -- 阿里云
channel（请在上面的渠道列表选择，输入序号）： 1
login as xxxxx success !
```

登录成功后 使用 matchvs debug <你生成gameServver 的 GS_key> 命令打开本地调试，信息如下：

```
matchvs debug 260e904d2daac445300592dc39395e45
        ==================== Develop config ====================
        SvcName:        svc-123456-0-0
        PodName:        deploy-123456-0-0-855667b-jkxmj
        RemoteHost:     directory10.matchvs.com
        RemotePort:     9982

2018/11/30 17:12:15 [I] [proxy_manager.go:298] proxy removed: []
2018/11/30 17:12:15 [I] [proxy_manager.go:308] proxy added: [matchvs]
2018/11/30 17:12:15 [I] [proxy_manager.go:331] visitor removed: []
2018/11/30 17:12:15 [I] [proxy_manager.go:340] visitor added: []
2018/11/30 17:12:15 [I] [control.go:240] [2ed21a3204021cee] login to server success, get run id [2ed21a3204021cee], server udp port [0]
2018/11/30 17:12:15 [I] [control.go:165] [2ed21a3204021cee] [matchvs] start proxy success
```

看到  start proxy success 消息说明本地调试打开成功。

执行 matchvs debug 成功后，会看到如下信息，这几条数据需要配置到 conf/config.toml 文件中。

```
==================== Develop config ====================
        SvcName:        svc-123456-0-0
        PodName:        deploy-123456-0-0-855667b-jkxmj
        RemoteHost:     directory10.matchvs.com
        RemotePort:     9982
```

SvcName, PodName，RemoteHost, RemotePort 与conf/config.toml 文件中的字段对应。

可参考如下修改：

```toml
# 独立部署参数配置
[Register]
Enable = false
GameID = 123456
# SvcName debug 调试中的 SvcName
SvcName = "svc-123456-0-0"
# PodName debug 调试中的 PodName
PodName = "deploy-123456-0-0-855667b-jkxmj"
# RemoteHost debug 调试中的 RemoteHost
RemoteHost = "directory10.matchvs.com"
# RemotePort debug 调试中的 RemotePort
RemotePort = 9982
LocalHost = "192.168.8.2"
LocalPort = 30054

# RoomManage 房间管理参数配置
[RoomManage]
# 需要使用房间管理 就设置为 true
Enable = false
# SvcName debug 调试中的 SvcName
SvcName = "svc-123456-0-0"
# PodName debug 调试中的 PodName
PodName = "deploy-123456-0-0-855667b-jkxmj"
# RemoteHost debug 调试中的 RemoteHost
RemoteHost = "directory10.matchvs.com"
# RemotePort debug 调试中的 RemotePort
RemotePort = 9982
```

#### 3、运行本地 gameServer 程序

假设在 windows 平台，在你的 golang 版本 gameServer 项目下面运行程序，比如我的工程是 gogs-demo ,目录结构如下：

```
┌ app/  		# 业务代码
├ conf/ 		# 配置文件
├ main.go 		# 程序入口
├ Dockerfile 	# docker部署，这个文件不要修改
├ Makefile 		# make 文件内容不能修改
└ gsmeta 		# 数据源文件内容不能修改
```

在 gogs-demo 目录下执行如下命令可以看到运行结果：

```verilog
go run main.go
2018/11/30 17:35:01.21 game_service.go:55 ▶ INFO 001 log level to set [DEBUG]
2018/11/30 17:35:01.21 room_manager.go:51 ▶ INFO 002 room manager is disable
2018/11/30 17:35:01.21 dir_register.go:46 ▶ INFO 003 register is disable
2018/11/30 17:35:01.21 stream_server.go:66 ▶ INFO 004 game server start success [0.0.0.0:30523]
```

我们在 app/app.go 文件中 OnCreateRoom 和 OnJoinRoom 等业务处理函数都有输出日志，当有人加入房间会输出如下日志：

```verilog
2018/11/30 17:37:31.02 app.go:33 ▶ DEBUG 005  OnCreateRoom map[mode:0 canWatch:2 createFlag:1 roomProperty: maxPlayer:3 state:1 createTime:1543570653 userID:1624685 userProfile:userProfile roomID:1708522620665729120]
2018/11/30 17:37:31.19 app.go:39 ▶ DEBUG 006  OnJoinRoom map[maxPlayers:3 checkins:[1624685] players:[1624685] roomID:1708522620665729120 userID:1624685 gameID:200773 userProfile:userProfile joinType:3]
```

但有人离开房间会输出如下日志：

```verilog
2018/11/30 17:38:51.37 app.go:104 ▶ DEBUG 008  OnDeleteRoom map[gameID:200773 roomID:1708522620665729120]
2018/11/30 17:38:51.41 app.go:57 ▶ DEBUG 00a  OnLeaveRoom map[gameID:200773 userID:1624685 roomID:1708522620665729120]
```

> **注意：** 这里只是简单的描述了接口触发，然后以日志的形式处理接口消息，相关的游戏逻辑由开发者自己去处理。

#### 4、退出本地调试

使用 matchvs debug 命令打开了本地调试后，如果不适用本地调试了并且没有使用命令关闭本地调试可能会出现，在调试客户端调用 `joinRoom` 相关接口的时候 出现 520 以及其他错误。

关闭本地调试需要在 开启 matchvs debug 后，使用 quit 命令退出。

```verilog
2018/11/30 17:12:15 [I] [proxy_manager.go:298] proxy removed: []
2018/11/30 17:12:15 [I] [proxy_manager.go:308] proxy added: [matchvs]
2018/11/30 17:12:15 [I] [proxy_manager.go:331] visitor removed: []
2018/11/30 17:12:15 [I] [proxy_manager.go:340] visitor added: []
2018/11/30 17:12:15 [I] [control.go:240] [2ed21a3204021cee] login to server success, get run id [2ed21a3204021cee], server udp port [0]
2018/11/30 17:12:15 [I] [control.go:165] [2ed21a3204021cee] [matchvs] start proxy success
......

quit
2018/11/30 17:47:17 [E] [control.go:148] [2ed21a3204021cee] work connection closed, broken pipe
2018/11/30 17:47:17 [W] [control.go:281] [2ed21a3204021cee] read error: broken pipe
2018/11/30 17:47:17 [I] [control.go:301] [2ed21a3204021cee] control writer is closing
```



## gameServer 发布

在发布上线之前，请确认你的游戏状态是处于 **已商用** 状态 [前往查看状态](http://www.matchvs.com/manage/gameContentList) 。 

发布后的环境，相对客户端SDK来说也就是release环境，所以 `login` 接口参数应该使用 `release` 

在编译可执行程序之前一定要把 example 目录下的 Makefile 文件拷贝到你的工程 root目录下。

#### 1、编译工程生成 linux 平台的可执行文件

编译可执行文件分为两种方法

- 使用 make 命令编译：在 windows 下需要安装 make 工具。[windows make 工具下载](http://gnuwin32.sourceforge.net/packages/make.htm) 。

- 使用 go 命令编译：需要指定编译环境。

第一种：使用 make 编译可执行文件会在你的项目目录下生成 docker_build 文件夹，内容结构如下：

```
├ docker_build/ 		#docker 打包文件 文件名不能修改
	├ conf/
		config.toml		# 配置文件
	├ Dockerfile		# docker 文件
	├ gameserver_go	 	# linux下的可执行程序 必须为 gameserver_go

```

第二种：使用 go 命令编译必须制定编译参数，生成的可执行文件名称必须为 gameserver_go，使用如下命令编译：

```powershell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gameserver_go
```

编程成功会在当前目录下面生成一个 gameserver_go 文件。然后手动新建 docker_build 文件夹，并且把 gameserver_go 、Dockerfile 、conf 文件夹 内容拷贝到 docker_build 文件夹目录下面。目录结构如上描述的。

#### 2、上传文件到仓库

编译好文件后，并且严格要求 docker_build 文件夹中的内容存在。使用 git 工具把工程代码提交到 matchvs 的gameServer 仓库中。提交同时要包含源代码。

> **注意** ：上传代码时一定要有 docker_build 文件和内容存在。并且文件中内容名称不能为其他的。

#### 3、官网发布

上传了编译好的工程文件，在官网对应的 gameServer 点击发布按钮。

发布成功后旁边的启动按钮就可以使用了，点击启动等待 gameServer 程序启动。启动成功后可以进入gameServer 中查看相关的日志啦。看到如下日志说明启动成功

```verilog
2018/11/30 17:35:01.21 game_service.go:55 ▶ INFO 001 log level to set [DEBUG]
2018/11/30 17:35:01.21 room_manager.go:51 ▶ INFO 002 room manager is disable
2018/11/30 17:35:01.21 dir_register.go:46 ▶ INFO 003 register is disable
2018/11/30 17:35:01.21 stream_server.go:66 ▶ INFO 004 game server start success [0.0.0.0:xxxxx]
```

#### 提示：

发布代码时一定要需要的目录结构和内容如下：

````
┌ docker_build/ 		#docker 打包文件 文件名不能修改
	├ conf/
		config.toml		# 配置文件
	├ Dockerfile		# docker 文件
	├ gameserver_go	 	# linux下的可执行程序 必须为 gameserver_go
├ Makefile  # make 文件
└ gsmeta 		# 数据源文件内容不能修改
````

但是为了你的代码管理最好把你工程的所有源代码也一起提交到你的 gameServer  git 仓库。



# Matchvs gameServer-go 框架使用

阅读这些内容之前，请先阅读 gameServer-go 入门文档。

gameServer-go 框架代码在src中，示例代码在 example 中。我们以 example 中的代码为示例对接口信息说明。

## 框架入口

在 example/main.go 文件中可看到如下代码：

```go
func main() {
	// handler 实现接口 matchvs.BaseGsHandler
	handler := app.NewApp()
	// 创建 gameServer
	gsserver := matchvs.NewGameServer(handler, "")
    // 设置消息推送
	handler.SetPushHandler(gsserver.GetPushHandler())
	// 启动 gameSever 服务
	go gsserver.Start()
	......
}
```

- NewApp()： 是 example/app/app.go 中的 App 类型。这个类型实现了接口  BaseGsHandler。

- NewGameServer()： 是 "github.com/matchvs/gameServer-go" 中的类型。使用 NewGameServer(...) 创建 gameServer 服务需要传入 BaseGsHandler 接口类型来处理相关的业务。
- SetPushHandler()：是设置消息推送业务类，这个类型是从 gameServer 类型中使用 GetPushHandler() 获取
- Start()：启动 gameServer 服务。

## 业务处理类型

定义业务处理类型要满足一下几个条件：

- 必须实现接口 BaseGsHandler ： 接受服务消息接口
- 必须定义 PushManager 类型成员：主动推送消息接口。

例如 example/app/app.go 的 App 类型定义

````go
type App struct {
    ......
    // push 用于主动推送消息的服务。
	push    matchvs.PushHandler
}
func (self *App) SetPushHandler(push matchvs.PushHandler) {
	self.push = push
}
......
````

所有的业务操作入口和出口都是在 App 这个业务类型中。开发者也可以自己修改定义这个类型。只要满足上面的两个条件即可。

## 数据接收业务接口

BaseGsHandler 接口类型是用于接收服务推送的消息。接口定义如下：

```go
type BaseGsHandler interface {
    // 消息接收
	message.IHandler
    // 消息推送
	SetPushHandler(PushHandler)
}
```

BaseGsHandler 扩展 IHandler 类型

```go
type IHandler interface {
	// 创建房间回调
	OnCreateRoom(*defines.MsOnCreateRoom) error
	// 加入房间回调
	OnJoinRoom(*defines.MsOnJoinRoom) error
	// 关闭房间回调
	OnJoinOver(map[string]interface{}) error
	// 打开房间回调
	OnJoinOpen(map[string]interface{}) error
	// 离开房间回调
	OnLeaveRoom(map[string]interface{}) error
	// 踢人回调
	OnKickPlayer(map[string]interface{}) error
	// 连接状态回调
	OnUserState(map[string]interface{}) error
	// 获取房间信息回调
	OnRoomDetail(*defines.MsRoomDetail) error
	// 设置房间属性回调
	OnSetRoomProperty(map[string]interface{}) error
	// 房间连接回调
	OnHotelConnect(map[string]interface{}) error
	// 消息广播
	OnReceiveEvent(*defines.MsOnReciveEvent) error
	// 房间断开
	OnDeleteRoom(map[string]interface{}) error
	// 连接房间检测回调
	OnHotelCheckin(map[string]interface{}) error
	// 设置帧同步
	OnSetFrameSyncRate(*defines.MsFrameSyncRateNotify) error
	// 帧数据更新
	OnFrameUpdate(*defines.FrameDataList) error
}
```

`IHandler` 接口中的函数是与 `Matchvs SDK` 中的操作对应，同时也与 `PushHandler`  主动推送接口中的操作对应。

例如：在 SDK 中有人触发了 `createRoom` 或者 `joinRoom` 接口，`gameServer` 会触发 `IHandler` 的 `OnCreateRoom` 和 `OnJoinRoom` 函数。



## 数据推送接口

PushHandler 接口用于 gameServer 主动推送数据。接口定义如下：

```go
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
	SetFrameSyncRate(gameID, frameRate, enableGS uint32, roomID uint64)
	// FrameBroadcast 发送帧同步数据给 游戏房间服务
	// gameID : 游戏ID
	// operation : 数据处理方式 0：只发客户端；1：只发GS；2：同时发送客户端和GS
	// roomID : 房间ID
	// cpProto : 要发送的数据
	FrameBroadcast(gameID uint32, operation int32, roomID uint64, cpProto []byte)
}
```











