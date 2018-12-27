package test

import (
	"encoding/json"
	"strconv"

	matchvs "github.com/matchvs/gameServer-go"
	"github.com/matchvs/gameServer-go/src/defines"
	"github.com/matchvs/gameServer-go/src/log"
)

type GsDefaultHandler struct {
	counter uint32
	push    matchvs.PushHandler
}

func (self *GsDefaultHandler) SetPushHandler(push matchvs.PushHandler) {
	self.push = push
}

// 创建房间回调
func (d *GsDefaultHandler) OnCreateRoom(req *defines.MsOnCreateRoom) (err error) {
	log.LogD(" OnCreateRoom %v", req)
	return
}

// 加入房间回调
func (d *GsDefaultHandler) OnJoinRoom(req *defines.MsOnJoinRoom) (err error) {
	log.LogD(" OnJoinRoom %v", req)
	return
}

// 关闭房间回调
func (d *GsDefaultHandler) OnJoinOver(req map[string]interface{}) (err error) {
	log.LogD(" OnJoinOver %v", req)
	return
}

// 打开房间回调
func (d *GsDefaultHandler) OnJoinOpen(req map[string]interface{}) (err error) {
	log.LogD(" OnJoinOpen %v", req)
	return
}

// 离开房间回调
func (d *GsDefaultHandler) OnLeaveRoom(req map[string]interface{}) (err error) {
	log.LogD(" OnLeaveRoom %v", req)
	return
}

// 踢人回调
func (d *GsDefaultHandler) OnKickPlayer(req map[string]interface{}) (err error) {
	log.LogD(" OnKickPlayer %v", req)
	return
}

// 连接状态回调
func (d *GsDefaultHandler) OnUserState(req map[string]interface{}) (err error) {
	log.LogD(" OnUserState %v", req)
	return
}

// 获取房间信息回调
func (d *GsDefaultHandler) OnRoomDetail(req *defines.MsRoomDetail) (err error) {
	log.LogD("OnRoomDetail %v", req)
	for _, v := range req.PlayersList {
		log.LogD("OnRoomDetail PlayersList %v", v)
	}
	log.LogD("OnRoomDetail WatchRoom %v", req.WatchRoom)
	return
}

// 设置房间属性回调
func (d *GsDefaultHandler) OnSetRoomProperty(req map[string]interface{}) (err error) {
	log.LogD(" OnSetRoomProperty %v", req)
	return
}

// 房间连接回调
func (d *GsDefaultHandler) OnHotelConnect(req map[string]interface{}) (err error) {
	log.LogD(" OnHotelConnect %v", req)
	return
}

// 消息广播
func (d *GsDefaultHandler) OnReceiveEvent(req *defines.MsOnReciveEvent) (err error) {
	// log.LogD(" OnReceiveEvent %v", string(req.CpProto))
	d.Example_Push(req)
	return
}

// 房间断开
func (d *GsDefaultHandler) OnDeleteRoom(req map[string]interface{}) (err error) {
	log.LogD(" OnHotelCloseConnect %v", req)
	return
}

// 连接房间检测回调
func (d *GsDefaultHandler) OnHotelCheckin(req map[string]interface{}) (err error) {
	log.LogD(" OnHotelCheckin %v", req)
	return
}

// 设置帧同步
func (d *GsDefaultHandler) OnSetFrameSyncRate(req *defines.MsFrameSyncRateNotify) (err error) {
	log.LogD(" OnHotelSetFrameSyncRate %v", req)
	return
}

// 帧数据更新
func (d *GsDefaultHandler) OnFrameUpdate(req *defines.MsFrameDataList) (err error) {
	for _, v := range req.Items {
		log.LogD(" OnFrameUpdate  len [%d]  CpProto [%s], Timestamp [%d]", len(req.Items), string(v.CpProto), v.Timestamp)
	}
	return
}

func (d *GsDefaultHandler) Example_Push(req *defines.MsOnReciveEvent) {
	var optMap map[string]interface{}
	if err := json.Unmarshal(req.CpProto, &optMap); err != nil {
		log.LogE("event message Unmarshal error %v", err)
		return
	}
	// log.LogD("event message [%v]", optMap)
	cmd := optMap["cmd"].(string)
	switch cmd {
	case "createRoom":
		d.example_createRoom(req.GameID)
	case "destroyRoom":
		roomID, _ := strconv.ParseUint(optMap["roomID"].(string), 10, 64)
		status, err := d.push.DestroyRoom(req.GameID, roomID)
		if err != nil {
			log.LogW("destroyRoom error %v", err)
			return
		}
		log.LogD("destroyRoom status [%d]", status)
	case "touchRoom":
		roomID, _ := strconv.ParseUint(optMap["roomID"].(string), 10, 64)
		ttl := optMap["ttl"].(float64)
		status, err := d.push.TouchRoom(req.GameID, uint32(ttl), roomID)
		if err != nil {
			log.LogW("touchRoom error %v", err)
			return
		}
		log.LogD("touchRoom status [%d]", status)
	case "end":
		log.LogD("收到消息数量：", d.counter)
	case "event":
		d.counter++
	case "kickPlayer":
		userID := optMap["userID"].(float64)
		d.push.KickPlayer(uint32(userID), req.RoomID)
	case "joinOver":
		d.push.JoinOpen(req.GameID, req.RoomID)
	case "joinOpen":
		d.push.JoinOver(req.GameID, req.RoomID)
	case "getRoomDetail":
		d.push.GetRoomDetail(req.GameID, 5, req.RoomID)
	case "setRoomProperty":
		d.push.SetRoomProperty(req.GameID, req.RoomID, "gameServer_go set room Property")
	case "setFrameSyncRate":
		setinfo := &defines.MsSetFrameSyncRateReq{}
		setinfo.FrameRate = uint32(optMap["frameRate"].(float64))
		setinfo.EnableGS = uint32(optMap["enableGS"].(float64))
		setinfo.RoomID = req.RoomID
		setinfo.GameID = req.GameID
		setinfo.CacheFrameMS = int32(optMap["cacheMs"].(float64))
		d.push.SetFrameSyncRate(setinfo)
	case "frameUpdate":
		data := []byte("test message from gameServer_go frame synchronization")
		operation := optMap["enableGS"].(float64)
		d.push.FrameBroadcast(req.GameID, int32(operation), req.RoomID, data)
	case "pushEvent":
		event := &defines.MsPushEventReq{
			PushType:  1,
			GameID:    req.GameID,
			RoomID:    req.RoomID,
			DestsList: req.DestsList[:],
			CpProto:   []byte("gameServer push event test golang"),
		}
		d.push.PushEvent(event)
	case "getCacheData":
		cacheMs := optMap["cacheMs"].(float64)
		d.push.GetOffLineCacheData(req.GameID, req.RoomID, int32(cacheMs))
	default:
	}
}

// 创建房间 示例
func (d *GsDefaultHandler) example_createRoom(gameID uint32) {
	crtm := new(defines.MsCreateRoomReq)
	crtm.GameID = gameID
	crtm.Ttl = 600
	crtm.RoomInfo = &defines.MsRoomInfo{
		RoomName:     "gameServer_go",
		MaxPlayer:    3,
		Mode:         1,
		CanWatch:     1,
		Visibility:   1,
		RoomProperty: "roomProperty gameServer_go",
	}
	crtm.WatchSet = &defines.MsWatchSeting{
		MaxWatch:        1,
		WatchPersistent: false,
		WatchDelayMs:    3000,
		CacheTime:       10000,
	}
	status, err := d.push.CreateRoom(crtm)
	if err != nil {
		log.LogW("createRoom error %v", err)
		return
	}
	log.LogD("createRoom Status [%d]", status)
}
