# Matchvs Golang 版 gameServer

使用 golang 版本的 gameServer 需要在本地安装好 golang 运行环境。如果还没有配置好 golang 环境的请先配置好。以下内容教程是默认你已经配置好并且熟悉 golang 开发环境。

## 新建 gameServer

前往 [matchvs 官网](http://www.matchvs.com/manage/gameServer) 新建一个 gameServer 服务，新建方法可以参考 [nodejs 版gameServer](https://doc.matchvs.com/QuickStart/GameServer-JavaScript) 版本的说明。



## 配置 gameServer 工程

#### 1、获取 matchvs golang 版本 gameServer 服务代码

使用如下命令获取拉取代码：

```shell
go get -u github.com/matchvs/gameServer-go
```

拉取成功后可以在你的 $GOPATH/github.com 目录下面看到 matchvs/gameServer-go 的服务代码。

在你的 `$GOPATH/github.com/ matchvs/gameServer-go` 目录下面有一个 example 文件，文件中的内容结构如下：

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

- 复制 `$GOPATH/github.com/ matchvs/gameServer-go/example/` 文件里面的内容到你刚刚拉下的仓库目录。

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

- 修改和 main.go 文件 import 中的内容。在 app.go 和 main.go import 文件中我们看到如下内容:

```go
"github.com/matchvs/gameServer-go/example/app"
```

  这条引用需要改为你自己项目的引用路径，比如你现在的项目名称为 gogs-demo ，那么你需要修改为：

```go
"gogs-demo/app"
```

  > 注意：以上内容是 golang 包导入的知识，如果不太理解请阅读 golang 相关的知识。
  >
  >
  >  如果没有修改 "github.com/matchvs/gameServer-go/example/app"  内容为你自己的项目内容，则项目运行将运行 example 中的 app.go

- 修改好 import 内容后就可以运行 第一个 gameServer 程序啦。运行内容大概如下：

```shell
  2018/11/30 16:26:43.63 game_service.go:55 ▶ INFO 001 log level to set [DEBUG]
  2018/11/30 16:26:43.63 room_manager.go:51 ▶ INFO 002 room manager is disable
  2018/11/30 16:26:43.63 dir_register.go:46 ▶ INFO 003 register is disable
  2018/11/30 16:26:43.63 stream_server.go:66 ▶ INFO 004 game server start success [0.0.0.0:36520]
```

####  [下一步 ]()   go gameServer 本地调试

文档正在努力输出中

#### [下一步 ]()   go gameServer 接口编程说明 

文档正在努力输出中