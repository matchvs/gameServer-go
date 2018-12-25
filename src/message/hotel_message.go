package message

import (
	pb "commonlibs/proto"
	"commonlibs/servers"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/matchvs/gameServer-go/src/defines"
	"github.com/matchvs/gameServer-go/src/log"
)

// 每个房间中的帧数据
type roomFrameData struct {
	FrameData map[uint32]*defines.MsFrameDataList
}

// 房间中帧数据缓存池
type roomFrameDataPool struct {
	cache map[uint64]*roomFrameData
	lock  sync.Mutex
}

func NewRoomFrameDataPool() *roomFrameDataPool {
	return &roomFrameDataPool{
		cache: make(map[uint64]*roomFrameData),
	}
}

func (self *roomFrameDataPool) addFrameData(gameID, frameIndex uint32, roomID uint64, item *defines.MsFrameDataItem) {
	var (
		roomframe *roomFrameData
		frameData *defines.MsFrameDataList
		ok        bool
	)
	self.lock.Lock()
	defer self.lock.Unlock()
	roomframe, ok = self.cache[roomID]

	if !ok {
		roomframe = new(roomFrameData)
		roomframe.FrameData = make(map[uint32]*defines.MsFrameDataList)
		self.cache[roomID] = roomframe
	}
	frameData, ok = roomframe.FrameData[frameIndex]
	if !ok {
		frameData = &defines.MsFrameDataList{}
		frameData.RoomID = roomID
		frameData.GameID = gameID
		frameData.FrameIndex = frameIndex
		frameData.Items = make([]*defines.MsFrameDataItem, 0, 100)
		frameData.Items = append(frameData.Items, item)
		roomframe.FrameData[frameIndex] = frameData
	} else {
		frameData.Items = append(frameData.Items, item)
	}
}

// 获取 或添加房间帧数据
func (self *roomFrameDataPool) getFrameData(gameID, frameIndex uint32, roomID uint64) *defines.MsFrameDataList {
	var (
		roomframe *roomFrameData
		frameData *defines.MsFrameDataList
		ok        bool
		noFrame   bool
	)
	self.lock.Lock()
	defer self.lock.Unlock()
	roomframe, ok = self.cache[roomID]
	if !ok {
		noFrame = true
	} else {
		frameData, ok = roomframe.FrameData[frameIndex]
		if !ok {
			noFrame = true
		}
	}
	if noFrame {
		frameData = &defines.MsFrameDataList{}
		frameData.RoomID = roomID
		frameData.GameID = gameID
		frameData.FrameIndex = frameIndex
		frameData.Items = make([]*defines.MsFrameDataItem, 0, 100)
	}
	return frameData
}

// 删除房间帧数据
func (self *roomFrameDataPool) delFrameData(frameIndex uint32, roomID uint64) {
	self.lock.Lock()
	roomframe, ok := self.cache[roomID]
	if ok {
		delete(roomframe.FrameData, frameIndex)
	}
	self.lock.Unlock()
}

func (self *roomFrameDataPool) delRoomFrame(roomID uint64) {
	self.lock.Lock()
	delete(self.cache, roomID)
	// log.LogD("delRoomFrame room frameSync number [%v]", len(self.cache))
	self.lock.Unlock()
}

type hotelRouters func(connID uint64, req servers.GSRequest) ([]byte, error)

type HotelMessage struct {
	MessageModel
	router            map[uint32]hotelRouters //命令路由
	clients           map[uint64]uint64
	enableFrameSync   map[uint64]bool
	enableGsFrameSync map[uint64]uint32
	roomFramesPool    *roomFrameDataPool
	lock              sync.RWMutex
}

func NewHotelModel(hd IHandler, cache *MessageCache) (ht *HotelMessage) {
	ht = new(HotelMessage)
	ht.router = make(map[uint32]hotelRouters)
	ht.clients = make(map[uint64]uint64)
	ht.handle = hd
	ht.msgCache = cache
	ht.enableFrameSync = make(map[uint64]bool)
	ht.enableGsFrameSync = make(map[uint64]uint32)
	ht.roomFramesPool = NewRoomFrameDataPool()
	ht.setRoute()
	return
}

// 设置路由
func (h *HotelMessage) setRoute() {
	if h.router == nil {
		h.router = make(map[uint32]hotelRouters)
	}
	//
	h.router[uint32(pb.HotelGsCmdID_HotelCreateConnect)] = h.onCreateConnect
	//
	h.router[uint32(pb.HotelGsCmdID_HotelBroadcastCMDID)] = h.onReceiveEvent
	//
	h.router[uint32(pb.HotelGsCmdID_HotelCloseConnet)] = h.onDeleteRoom
	//
	h.router[uint32(pb.HotelGsCmdID_HotelPlayerCheckin)] = h.onPlayerCheckin

	h.router[uint32(pb.HotelGsCmdID_GSSetFrameSyncRateNotifyCMDID)] = h.onSetFrameSyncRate

	h.router[uint32(pb.HotelGsCmdID_GSFrameDataNotifyCMDID)] = h.onFrameDataNotify

	h.router[uint32(pb.HotelGsCmdID_GSFrameSyncNotifyCMDID)] = h.onFrameUpdate
}

// 判断是不是 hotel 处理的模块
func (m *HotelMessage) CanDeal(cmdid int32) bool {
	_, ok := pb.HotelGsCmdID_name[cmdid]
	return ok
}

// 添加链接
func (h *HotelMessage) addClient(roomID, userid uint64) {
	if h.clients == nil {
		h.clients = make(map[uint64]uint64)
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	h.clients[roomID] = userid
}

// 获取链接
func (h *HotelMessage) GetClient(roomID uint64) (id uint64, ok bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	id, ok = h.clients[roomID]
	return
}

// 删除链接
func (h *HotelMessage) DelClient(roomID uint64) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.clients != nil {
		delete(h.clients, roomID)
	}
}

func (h *HotelMessage) delFrameSycnCache(roomID uint64) {
	delete(h.enableFrameSync, roomID)
	delete(h.enableGsFrameSync, roomID)
	h.roomFramesPool.delRoomFrame(roomID)
}

// 消息路由处理
// req 收到 hotel 推送的消息
func (h *HotelMessage) Route(connID uint64, req servers.GSRequest) (res []byte, err error) {
	//log.LogD("hotel route cmdid [%d]", req.CmdId)
	// 获取业务处理函数

	handler, ok := h.router[req.CmdId]
	if !ok {
		log.LogW("hotel no this cmdid [%d]", req.CmdId)
		return
	}

	res, err = handler(connID, req)

	if err != nil {
		log.LogE("handler error %v", err)
	}
	return
}

func (m *HotelMessage) onCreateConnect(connID uint64, req servers.GSRequest) ([]byte, error) {
	connectInfo := &pb.Connect{}
	proto.Unmarshal(req.Message, connectInfo)
	// log.LogD("onCreateConnect [%v]", connectInfo)
	m.addClient(connectInfo.RoomID, connID)
	// 回复
	ack := &pb.ConnectAck{
		Status: uint32(pb.ErrorCode_OK),
	}
	return proto.Marshal(ack)
}

func (m *HotelMessage) onReceiveEvent(connID uint64, req servers.GSRequest) ([]byte, error) {
	broad := &pb.HotelBroadcast{}
	proto.Unmarshal(req.Message, broad)
	m.addClient(broad.RoomID, connID)

	event := &defines.MsOnReciveEvent{
		GameID:  broad.GameID,
		RoomID:  broad.RoomID,
		UserID:  broad.UserID,
		Flag:    broad.Flag,
		CpProto: broad.CpProto[:],
	}

	m.handle.OnReceiveEvent(event)

	ack := &pb.HotelBroadcastAck{
		UserID: uint32(req.UserId),
		Status: uint32(pb.ErrorCode_OK),
	}
	return proto.Marshal(ack)
}

func (m *HotelMessage) onDeleteRoom(connID uint64, req servers.GSRequest) ([]byte, error) {
	closeConn := &pb.CloseConnect{}
	proto.Unmarshal(req.Message, closeConn)
	// delete the client connect map
	m.DelClient(closeConn.RoomID)
	// clear all cache frame sysnc data
	m.delFrameSycnCache(closeConn.RoomID)
	delInfo := make(map[string]interface{})
	delInfo["gameID"] = closeConn.GetGameID()
	delInfo["roomID"] = closeConn.GetRoomID()
	// 业务处理
	m.handle.OnDeleteRoom(delInfo)

	ack := &pb.CloseConnectAck{
		Status: uint32(pb.ErrorCode_OK),
	}
	return proto.Marshal(ack)
}

func (m *HotelMessage) onPlayerCheckin(connID uint64, req servers.GSRequest) ([]byte, error) {

	playercheckin := &pb.PlayerCheckin{}
	proto.Unmarshal(req.Message, playercheckin)
	m.addClient(playercheckin.RoomID, connID)

	//加入房间的信息
	roominfo := &pb.JoinExtInfo{}
	joinroom := m.msgCache.GetWaitJoin(playercheckin.RoomID, playercheckin.UserID)
	if joinroom == nil {
		return nil, nil
	}
	proto.Unmarshal(joinroom.CpProto, roominfo)
	m.msgCache.DelWaitJoin(playercheckin.RoomID, playercheckin.UserID)

	onjoin := &defines.MsOnJoinRoom{
		RoomID:      playercheckin.RoomID,
		UserID:      playercheckin.UserID,
		GameID:      playercheckin.GameID,
		UserProfile: roominfo.GetUserProfile(),
		JoinType:    roominfo.GetJoinType(),
		MaxPlayers:  playercheckin.GetMaxPlayers(),
		Checkins:    playercheckin.GetCheckins(),
		Players:     playercheckin.GetPlayers(),
	}

	// 业务处理
	m.handle.OnJoinRoom(onjoin)
	ack := &pb.PlayerCheckinAck{
		Status: uint32(pb.ErrorCode_OK),
	}
	return proto.Marshal(ack)
}

func (m *HotelMessage) onSetFrameSyncRate(connID uint64, req servers.GSRequest) ([]byte, error) {
	syncRate := &pb.GSSetFrameSyncRateNotify{}
	proto.Unmarshal(req.Message, syncRate)
	m.addClient(syncRate.RoomID, connID)

	if syncRate.FrameRate == 0 {
		m.enableFrameSync[syncRate.RoomID] = false
		m.roomFramesPool.delRoomFrame(syncRate.RoomID)
	} else {
		m.enableFrameSync[syncRate.RoomID] = true
	}
	frameSyncRate := &defines.MsFrameSyncRateNotify{
		FrameIdx:     syncRate.GetFrameIdx(),
		FrameRate:    syncRate.GetFrameRate(),
		RoomID:       syncRate.GetRoomID(),
		EnableGS:     syncRate.GetEnableGS(),
		Timestamp:    syncRate.GetTimeStamp(),
		GameID:       syncRate.GetGameID(),
		CacheFrameMS: syncRate.GetCacheFrameMS(),
	}
	m.enableGsFrameSync[syncRate.RoomID] = syncRate.EnableGS
	m.handle.OnSetFrameSyncRate(frameSyncRate)
	return nil, nil
}

func (m *HotelMessage) onFrameDataNotify(connID uint64, req servers.GSRequest) ([]byte, error) {
	dataNotify := &pb.GSFrameDataNotify{}
	proto.Unmarshal(req.Message, dataNotify)
	m.addClient(dataNotify.RoomID, connID)

	if ok := m.enableFrameSync[dataNotify.RoomID]; ok {
		frameItem := &defines.MsFrameDataItem{
			SrcUserID: dataNotify.SrcUid,
			CpProto:   dataNotify.CpProto[:],
			Timestamp: dataNotify.TimeStamp,
		}
		m.roomFramesPool.addFrameData(dataNotify.GameID, dataNotify.FrameIdx, dataNotify.RoomID, frameItem)
	} else {
		m.roomFramesPool.delRoomFrame(dataNotify.RoomID)
	}
	return nil, nil
}

func (m *HotelMessage) onFrameUpdate(connID uint64, req servers.GSRequest) ([]byte, error) {
	frameUpdate := &pb.GSFrameSyncNotify{}
	proto.Unmarshal(req.Message, frameUpdate)
	m.addClient(frameUpdate.RoomID, connID)
	if ok := m.enableFrameSync[frameUpdate.RoomID]; ok {
		frameData := m.roomFramesPool.getFrameData(frameUpdate.GameID, frameUpdate.LastIdx, frameUpdate.RoomID)
		m.roomFramesPool.delFrameData(frameUpdate.LastIdx, frameUpdate.RoomID)
		m.handle.OnFrameUpdate(frameData)

	} else {
		m.roomFramesPool.delRoomFrame(frameUpdate.RoomID)
	}
	return nil, nil
}

func (m *HotelMessage) onGetCacheData(connID uint64, req servers.GSRequest) ([]byte, error) {
	log.LogD("触发获取断线帧数据")
	return nil, nil
}
