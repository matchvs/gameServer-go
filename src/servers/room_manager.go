package servers

import (
	"fmt"
	"hash/crc32"
	"os"
	"strconv"

	"github.com/matchvs/gameServer-go/src/config"
	"github.com/matchvs/gameServer-go/src/log"
	pb "github.com/matchvs/gameServer-go/src/proto"

	"github.com/golang/protobuf/proto"
)

type RoomManager struct {
	simple     *SimpleClient
	svcName    string
	podName    string
	remoteHost string
	remotePort int32
	reserved   uint64
	userID     uint32
	enable     bool
}

func NewRoomManager(conf *config.RoomConf) (r *RoomManager) {
	r = new(RoomManager)
	r.reserved = 0
	r.svcName = os.Getenv("DIRECTORY_SVC_NAME")
	if r.svcName != "" {
		r.podName = os.Getenv("DIRECTORY_POD_NAME")
		r.remoteHost = os.Getenv("DIRECTORY_REMOTE_HOST")
		port, _ := strconv.ParseInt(os.Getenv("DIRECTORY_REMOTE_PORT"), 10, 32)
		r.remotePort = int32(port)
	} else {
		r.svcName = conf.SvcName
		r.podName = conf.PodName
		r.remoteHost = conf.RemoteHost
		r.remotePort = conf.RemotePort
	}
	r.enable = conf.Enable
	r.userID = crc32.ChecksumIEEE([]byte(r.svcName + r.podName))
	addr := fmt.Sprintf("%s:%d", r.remoteHost, r.remotePort)
	r.simple = NewSimpleClient(addr)
	r.connect()
	return
}

func (r *RoomManager) connect() {
	if r.enable == false {
		log.LogW("room configuration is disable")
		return
	}
	err := r.simple.Run()
	if err != nil {
		log.LogE("room manage run error :%v", err)
		return
	}
	log.LogD("room manage run ...")
}

func (r *RoomManager) Stop() {
	r.simple.Close()
}

func (r *RoomManager) getReverse() uint32 {
	r.reserved++
	if r.reserved > maxReserved {
		r.reserved = 0
	}
	return uint32(r.reserved)
}

func (r *RoomManager) sendMsg(cmdID uint32, msg []byte) (*pb.Package_Frame, error) {
	req := &pb.Package_Frame{
		Type:     pb.Package_LeagueMessage,
		Version:  VERSION,
		CmdId:    cmdID,
		Reserved: r.getReverse(),
		UserId:   uint64(r.userID),
		Message:  msg[:len(msg)],
	}
	log.LogD("room manager send message [:%v]", req)
	resp, err := r.simple.Send(req)
	if err != nil {
		log.LogE("room manage send message error :%v", err)
		return nil, err
	}
	return resp, nil
}

// 创建房间
func (r *RoomManager) CreateRoom(ct *pb.CreateRoom) (*pb.CreateRoomAck, error) {
	ct.PodName = r.podName
	ct.SvcName = r.svcName
	log.LogD("CreateRoom request %v", *ct)
	msg, err := proto.Marshal(ct)

	if err != nil {
		log.LogE("create room marshal error :%v", err)
		return nil, err
	}
	// 请求
	resp, err := r.sendMsg(uint32(pb.GSDirectoryCmdID_GSCreateRoomCmd), msg)
	if err != nil {
		log.LogE("create room request error :%v", err)
		return nil, err
	}
	// 回复
	ack := new(pb.CreateRoomAck)
	if err = proto.Unmarshal(resp.Message, ack); err != nil {
		log.LogE("create room response ummarshal %v", err)
		return nil, err
	}
	return ack, nil
}

// 销毁房间
func (r *RoomManager) DestroyRoom(dt *pb.DestroyRoom) (*pb.DestroyRoomAck, error) {
	dt.PodName = r.podName
	dt.SvcName = r.svcName
	log.LogD("DestroyRoom request %v", *dt)
	msg, err := proto.Marshal(dt)
	if err != nil {
		log.LogE("destroy room marshal error :%v", err)
		return nil, err
	}
	resp, err := r.sendMsg(uint32(pb.GSDirectoryCmdID_GSDestroyRoomCmd), msg)
	if err != nil {
		log.LogE("destroy room request error :%v", err)
		return nil, err
	}
	ack := new(pb.DestroyRoomAck)
	if err = proto.Unmarshal(resp.Message, ack); err != nil {
		log.LogE("destroy room response ummarshal %v", err)
		return nil, err
	}
	return ack, nil
}

// 修改房间存活时长
func (r *RoomManager) TouchRoom(dt *pb.TouchRoom) (*pb.TouchRoomAck, error) {
	dt.PodName = r.podName
	dt.SvcName = r.svcName

	log.LogD("TouchRoom request %v", *dt)

	msg, err := proto.Marshal(dt)
	if err != nil {
		log.LogE("touch room marshal error :%v", err)
		return nil, err
	}
	resp, err := r.sendMsg(uint32(pb.GSDirectoryCmdID_GSTouchRoomCmd), msg)
	if err != nil {
		log.LogE("touch room request error :%v", err)
		return nil, err
	}
	ack := new(pb.TouchRoomAck)
	if err = proto.Unmarshal(resp.Message, ack); err != nil {
		log.LogE("touch room response ummarshal %v", err)
		return nil, err
	}
	return ack, nil
}
