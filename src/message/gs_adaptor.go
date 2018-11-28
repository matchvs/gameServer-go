package message

import (
	pb "commonlibs/proto"
	"commonlibs/servers"
	"errors"
	"sync"

	"github.com/matchvs/gameServer-go/src/log"
)

type GSAdaptor struct {
	mvsModel   *MvsMessage
	hotelModel *HotelMessage
	clients    map[uint64]servers.GSWrite
	lock       sync.RWMutex
	msgCache   *MessageCache
}

func NewGSAdaptor(h IHandler) (gs *GSAdaptor) {
	gs = new(GSAdaptor)
	gs.msgCache = NewMessageCache()
	gs.mvsModel = NewMvsModel(h, gs.msgCache)
	gs.hotelModel = NewHotelModel(h, gs.msgCache)
	gs.clients = make(map[uint64]servers.GSWrite)
	return
}

func (g *GSAdaptor) addClient(key uint64, value servers.GSWrite) {
	g.lock.Lock()
	if g.clients == nil {
		g.clients = make(map[uint64]servers.GSWrite)
	}
	g.clients[key] = value
	g.lock.Unlock()
}
func (g *GSAdaptor) delClient(key uint64) {
	g.lock.Lock()
	if len(g.clients) > 0 {
		delete(g.clients, key)
	}
	g.lock.Unlock()
	// g.clients.Delete(userid)
}
func (g *GSAdaptor) getClient(key uint64) (servers.GSWrite, bool) {
	g.lock.Lock()
	value, ok := g.clients[key]
	g.lock.Unlock()
	return value, ok
}

// 建立连接
func (g *GSAdaptor) OnConnect(userid uint64, token string, write servers.GSWrite) error {
	g.addClient(userid, write)
	return nil
}

// 断开连接
func (g *GSAdaptor) OnDisconnect(userid uint64, token string) error {
	g.mvsModel.DelClient(userid)
	return nil
}

// 收到客户端消息查找指定的模块处理
// req *pb.Package_Frame 类型, rsp 是 func(*pb.Package_Frame) error 类型
func (g *GSAdaptor) Route(connID uint64, req servers.GSRequest, write servers.GSWrite) (err error) {

	var (
		CmdHeartbeat = 99999999
		resData      []byte
	)
	// 心跳处理
	if req.CmdId == uint32(CmdHeartbeat) {
		write(req)
		// log.LogD("心跳 [%v]", req)
		return nil
	}

	if g.mvsModel.CanDeal(int32(req.CmdId)) {
		//查找 mvs 路由处理
		resData, err = g.mvsModel.Route(connID, req)
		if err != nil {
			log.LogE("mvs [%d] handler error %v", req.CmdId, err)
			return err
		}
	} else if g.hotelModel.CanDeal(int32(req.CmdId)) {
		//查找 hotel 路由处理
		resData, err = g.hotelModel.Route(connID, req)
		if err != nil {
			log.LogE("hotel [%d] handler error %v", req.CmdId, err)
			return err
		}
	} else {
		log.LogW("no router [%d]", req.CmdId)
		return
	}
	resp := &pb.Package_Frame{
		Type:     req.Type,
		CmdId:    req.CmdId + 1,
		Version:  servers.VERSION,
		UserId:   req.UserId,
		Reserved: req.Reserved,
		Message:  resData,
	}
	// log.LogD(" 回复包：%v", resp)
	if err := write(resp); err != nil {
		log.LogE("client [%d] ack error %v", req.CmdId, err)
		return err
	}
	return nil
}

func (g *GSAdaptor) push(userid uint64, cmdid uint32, message []byte) error {
	resp := &pb.Package_Frame{
		Type:    pb.Package_PushMessage,
		CmdId:   cmdid,
		Version: servers.VERSION,
		UserId:  userid,
		Message: message[:len(message)],
	}
	write, ok := g.getClient(userid)
	if !ok {
		log.LogW("no this client %d", userid)
		return nil
	}
	// value, _ := g.clients.Load(userid)
	// write := value.(GSWrite)
	if err := write(resp); err != nil {
		log.LogE("send client message error : %v", err)
		return err
	}
	return nil
}

// 发送消息
func (g *GSAdaptor) PushMvs(cmdid uint32, message []byte) error {
	userID := g.mvsModel.GetClient()
	g.push(userID, cmdid, message)
	return nil
}

func (g *GSAdaptor) PushHotel(cmdid uint32, roomID uint64, message []byte) error {
	ok := g.hotelModel.CanDeal(int32(cmdid))
	if !ok {
		log.LogW("no cmdID")
		return errors.New("no cmdID")
	}
	userID, ok := g.hotelModel.GetClient(roomID)
	if !ok {
		log.LogW("no room")
		return errors.New("no room")
	}
	g.push(userID, cmdid, message)
	return nil
}
