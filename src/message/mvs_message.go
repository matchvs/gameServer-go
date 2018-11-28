package message

import (
	"math/rand"
	"strconv"

	pb "github.com/matchvs/gameServer-go/src/proto"
	"github.com/matchvs/gameServer-go/src/servers"

	"github.com/matchvs/gameServer-go/src/defines"
	"github.com/matchvs/gameServer-go/src/log"

	"github.com/deckarep/golang-set"

	"github.com/golang/protobuf/proto"
)

// Mvs 消息路由类型
type mvsRouters func(*pb.Request) ([]byte, error)

type MvsMessage struct {
	MessageModel
	router  map[uint32]mvsRouters
	clients mapset.Set
}

// Mvs 通信处理模块
func NewMvsModel(h IHandler, cache *MessageCache) (mvs *MvsMessage) {
	mvs = new(MvsMessage)
	mvs.handle = h
	mvs.msgCache = cache
	mvs.router = make(map[uint32]mvsRouters)
	mvs.setRoute()
	mvs.clients = mapset.NewSet()
	return
}

func (m *MvsMessage) AddClient(id uint64) {
	if m.clients == nil {
		m.clients = mapset.NewSet()
	}
	m.clients.Add(id)
}

func (m *MvsMessage) DelClient(id uint64) {
	m.clients.Remove(id)
}

func (m *MvsMessage) GetClient() (id uint64) {
	arrclients := (m.clients.ToSlice())
	index := rand.Intn(len(arrclients))
	if index >= len(arrclients) {
		index = 0
	}
	id = (arrclients[index]).(uint64)
	return
}

// 设置路由
func (m *MvsMessage) setRoute() {
	if m.router == nil {
		m.router = make(map[uint32]mvsRouters)
	}
	// 创建房间
	m.router[uint32(pb.MvsGsCmdID_MvsCreateRoomReq)] = m.OnCreateRoom
	// 加入房间
	m.router[uint32(pb.MvsGsCmdID_MvsJoinRoomReq)] = m.OnJoinRoom
	// 关闭房间
	m.router[uint32(pb.MvsGsCmdID_MvsJoinOverReq)] = m.OnJoinOver
	// 打开房间
	m.router[uint32(pb.MvsGsCmdID_MvsJoinOpenReq)] = m.OnJoinOpen
	// 离开房间
	m.router[uint32(pb.MvsGsCmdID_MvsLeaveRoomReq)] = m.OnLeaveRoom

	m.router[uint32(pb.MvsGsCmdID_MvsKickPlayerReq)] = m.OnKickPlayer

	m.router[uint32(pb.MvsGsCmdID_MvsNetworkStateReq)] = m.OnUserState

	m.router[uint32(pb.MvsGsCmdID_MvsGetRoomDetailPush)] = m.OnRoomDetail

	m.router[uint32(pb.MvsGsCmdID_MvsSetRoomPropertyReq)] = m.OnSetRoomProperty
}

// 判断是不是 mvs 处理的模块
func (m *MvsMessage) CanDeal(cmdid int32) bool {
	_, ok := pb.MvsGsCmdID_name[cmdid]
	return ok
}

// 消息路由处理
func (m *MvsMessage) Route(connID uint64, req servers.GSRequest) (res []byte, err error) {
	// log.LogD("mvs route cmdid [%d]", req.CmdId)
	m.AddClient(connID)

	handler, ok := m.router[req.CmdId]
	if !ok {
		log.LogW("mvs no this cmdid [%d]", req.CmdId)
		return
	}

	request := &pb.Request{}
	if err = proto.Unmarshal(req.Message, request); err != nil {
		log.LogE("request message Unmarshal error:", err)
		return
	}

	_, err = handler(request)
	if err != nil {
		log.LogE("handler error %v", err)
	}

	reply := &pb.Reply{
		UserID: request.GetUserID(),
		GameID: request.GetGameID(),
		RoomID: request.GetRoomID(),
	}
	res, _ = proto.Marshal(reply)
	return
}

// 有人创建房间
func (m *MvsMessage) OnCreateRoom(req *pb.Request) ([]byte, error) {
	createinfo := &pb.CreateExtInfo{}
	proto.Unmarshal(req.CpProto, createinfo)
	infoMap := make(map[string]interface{})
	infoMap["userID"] = createinfo.UserID
	infoMap["userProfile"] = string(createinfo.UserProfile)
	infoMap["roomID"] = createinfo.RoomID
	infoMap["roomProperty"] = string(createinfo.RoomProperty)
	infoMap["maxPlayer"] = createinfo.MaxPlayer
	infoMap["state"] = createinfo.State
	infoMap["mode"] = createinfo.Mode
	infoMap["canWatch"] = createinfo.CanWatch
	infoMap["createFlag"] = createinfo.CreateFlag
	infoMap["createTime"] = createinfo.CreateTime
	m.handle.OnCreateRoom(infoMap)
	return nil, nil
}

func (m *MvsMessage) OnJoinRoom(req *pb.Request) ([]byte, error) {
	cache := &pb.Request{
		UserID:  req.UserID,
		GameID:  req.GameID,
		RoomID:  req.RoomID,
		CpProto: req.CpProto,
	}
	m.msgCache.AddWaitJoin(req.RoomID, req.UserID, cache)
	return nil, nil
}
func (m *MvsMessage) OnJoinOver(req *pb.Request) ([]byte, error) {

	infoMap := make(map[string]interface{})
	infoMap["gameID"] = req.GetGameID()
	infoMap["userID"] = req.GetUserID()
	infoMap["roomID"] = req.GetRoomID()
	m.handle.OnJoinOver(infoMap)
	return nil, nil
}
func (m *MvsMessage) OnJoinOpen(req *pb.Request) ([]byte, error) {

	infoMap := make(map[string]interface{})
	infoMap["gameID"] = req.GetGameID()
	infoMap["userID"] = req.GetUserID()
	infoMap["roomID"] = req.GetRoomID()
	m.handle.OnJoinOpen(infoMap)
	return nil, nil
}
func (m *MvsMessage) OnLeaveRoom(req *pb.Request) ([]byte, error) {
	infoMap := make(map[string]interface{})
	roomID := req.GetRoomID()
	userID := req.GetUserID()
	infoMap["gameID"] = req.GetGameID()
	infoMap["userID"] = userID
	infoMap["roomID"] = roomID
	m.handle.OnLeaveRoom(infoMap)
	m.msgCache.DelWaitJoin(roomID, userID)
	return nil, nil
}

func (m *MvsMessage) OnKickPlayer(req *pb.Request) ([]byte, error) {
	infoMap := make(map[string]interface{})
	roomID := req.GetRoomID()
	userID := req.GetUserID()
	infoMap["gameID"] = req.GetGameID()
	infoMap["userID"] = userID
	infoMap["roomID"] = roomID
	m.handle.OnKickPlayer(infoMap)
	m.msgCache.DelWaitJoin(roomID, userID)
	return nil, nil
}

func (m *MvsMessage) OnUserState(req *pb.Request) ([]byte, error) {
	infoMap := make(map[string]interface{})
	roomID := req.GetRoomID()
	userID := req.GetUserID()
	state, _ := strconv.ParseInt(string(req.CpProto), 10, 32)
	infoMap["gameID"] = req.GetGameID()
	infoMap["userID"] = userID
	infoMap["roomID"] = roomID
	infoMap["state"] = state
	m.handle.OnUserState(infoMap)
	log.LogD("OnUserState state:[%s]", string(req.CpProto))
	if state == 3 {
		m.msgCache.DelWaitJoin(roomID, userID)
	}

	return nil, nil
}

func (m *MvsMessage) OnSetRoomProperty(req *pb.Request) ([]byte, error) {
	infoMap := make(map[string]interface{})
	infoMap["gameID"] = req.GetGameID()
	infoMap["userID"] = req.GetUserID()
	infoMap["roomID"] = req.GetRoomID()
	m.handle.OnSetRoomProperty(infoMap)
	return nil, nil
}

func (m *MvsMessage) OnRoomDetail(req *pb.Request) ([]byte, error) {
	roomDetail := &pb.RoomDetail{}
	proto.Unmarshal(req.CpProto, roomDetail)
	room := &defines.MsRoomDetail{}

	room.RoomID = req.GetRoomID()
	room.GameID = req.GetGameID()
	room.State = roomDetail.GetState()
	room.CanWatch = roomDetail.GetCanWatch()
	room.CreateFlag = roomDetail.GetCreateFlag()
	room.MaxPlayer = roomDetail.GetMaxPlayer()
	room.Mode = roomDetail.GetMode()
	room.Owner = roomDetail.GetOwner()
	room.RoomProperty = string(roomDetail.GetRoomProperty())
	room.UserID = req.GetUserID()

	room.PlayersList = make([]*defines.MsPlayerInfo, len(roomDetail.PlayerInfos))
	for k, v := range roomDetail.PlayerInfos {
		player := new(defines.MsPlayerInfo)
		player.UserID = v.GetUserID()
		player.UserProfile = string(v.GetUserProfile())
		room.PlayersList[k] = player
	}

	watchRoom := new(defines.MsWatchRoom)
	pb_watchRoom := roomDetail.GetWatchRoom()
	watchRoom.State = pb_watchRoom.WatchInfo.GetState()
	watchRoom.RoomID = pb_watchRoom.WatchInfo.GetRoomID()
	watchRoom.CurWatch = pb_watchRoom.WatchInfo.GetCurWatch()

	watchSet := new(defines.MsWatchSeting)
	watchSet.MaxWatch = pb_watchRoom.WatchInfo.WatchSetting.GetMaxWatch()
	watchSet.CacheTime = pb_watchRoom.WatchInfo.WatchSetting.GetCacheTime()
	watchSet.WatchDelayMs = pb_watchRoom.WatchInfo.WatchSetting.GetWatchDelayMs()
	watchSet.WatchPersistent = pb_watchRoom.WatchInfo.WatchSetting.GetWatchPersistent()
	watchRoom.WatchSet = watchSet

	watchRoom.WatchPlayersList = make([]*defines.MsPlayerInfo, len(pb_watchRoom.WatchPlayers))
	for k, v := range pb_watchRoom.WatchPlayers {
		watchPlayer := new(defines.MsPlayerInfo)
		watchPlayer.UserProfile = string(v.GetUserProfile())
		watchPlayer.UserID = v.GetUserID()
		watchRoom.WatchPlayersList[k] = watchPlayer
	}

	room.WatchRoom = watchRoom
	m.handle.OnRoomDetail(room)
	return nil, nil
}
