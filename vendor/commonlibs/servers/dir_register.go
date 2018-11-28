package servers

import (
	pb "commonlibs/proto"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/matchvs/gameServer-go/src/config"
	"github.com/matchvs/gameServer-go/src/log"
)

var pushHandler *RegisterPushHandler

type RegisterPushHandler struct {
}

func (self *RegisterPushHandler) RecvPush(req *pb.Package_Frame) error {
	data, _ := proto.Marshal(req)
	log.LogD("收到消息：[CmdID=%d][UserID=%d][length=%v]", req.CmdId, req.UserId, len(data))
	return nil
}

type DirRegister struct {
	stream  *StreamClient
	regConf *config.RegisterConf
	stop    bool
	load    uint32
}

func NewRegister(conf *config.RegisterConf) (reg *DirRegister) {
	addr := fmt.Sprintf("%s:%d", conf.RemoteHost, conf.RemotePort)
	reg = new(DirRegister)
	pushHandler = new(RegisterPushHandler)
	reg.stream = NewStreamClient(addr, crc32.ChecksumIEEE([]byte(conf.SvcName+conf.PodName)), pushHandler)
	reg.stop = false
	reg.load = 0
	reg.regConf = conf
	return
}

func (r *DirRegister) Run() {

	if r.regConf.Enable == false {
		log.LogI("register is disable")
		return
	}

	if err := r.stream.Start(); err != nil {
		log.LogE("register fail remote addr : %s:%d", r.regConf.RemoteHost, r.regConf.RemotePort)
		return
	}
	log.LogI("register success : %s:%d", r.regConf.RemoteHost, r.regConf.RemotePort)
	r.Login()
	go func() {
		for {
			if r.stop {
				r.Loginout()
				r.stream.Stop()
				break
			}
			if r.stream.ConnectSuccess() == false {
				log.LogW("disconnect server : %s:%d", r.regConf.RemoteHost, r.regConf.RemotePort)
				break
			}
			r.heartBeat()
			time.Sleep(time.Second * time.Duration(5))
		}
	}()
}
func (r *DirRegister) Stop() {
	r.stop = true
}

func (r *DirRegister) Login() {
	pkg := &pb.GSLogin{}
	pkg.GameID = r.regConf.GameID
	pkg.SvcName = r.regConf.SvcName
	pkg.PodName = r.regConf.PodName
	pkg.Host = r.regConf.LocalHost
	pkg.Port = r.regConf.LocalPort
	msg, _ := proto.Marshal(pkg)
	r.stream.Push(uint32(pb.GSDirectoryCmdID_GSLoginCmd), msg)
}

func (r *DirRegister) Loginout() {
	pkg := &pb.GSLogout{}
	pkg.GameID = r.regConf.GameID
	pkg.SvcName = r.regConf.SvcName
	pkg.PodName = r.regConf.PodName
	msg, _ := proto.Marshal(pkg)
	r.stream.Push(uint32(pb.GSDirectoryCmdID_GSLogoutCmd), msg)
}

func (r *DirRegister) heartBeat() {
	hb := &pb.GSHeartbeat{
		Load: r.load,
	}
	message, _ := proto.Marshal(hb)
	r.stream.Push(uint32(pb.GSDirectoryCmdID_GSHeartbeatCmd), message)
}
