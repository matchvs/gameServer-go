/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-11-28 14:30:33
 * @LastEditors: Ville
 * @LastEditTime: 2018-11-28 17:15:56
 * @Description: game server handler module, the struct of  App  implemente the interface which is located in game_server.go
 				 it is named BaseInterface
*/

package app

import (
	"encoding/json"
	"strconv"

	"github.com/matchvs/gameServer-go"
	"github.com/matchvs/gameServer-go/src/defines"
	"github.com/matchvs/gameServer-go/src/log"
)

type App struct {
	counter uint32
	push    *matchvs.PushManager
}

func (self *App) SetPushHandler(push *matchvs.PushManager) {
	self.push = push
}

// 创建房间回调
func (d *App) OnCreateRoom(req map[string]interface{}) (err error) {
	log.LogD(" OnCreateRoom %v", req)
	return
}

// 加入房间回调
func (d *App) OnJoinRoom(req map[string]interface{}) (err error) {
	log.LogD(" OnJoinRoom %v", req)
	return
}

// 关闭房间回调
func (d *App) OnJoinOver(req map[string]interface{}) (err error) {
	log.LogD(" OnJoinOver %v", req)
	return
}

// 打开房间回调
func (d *App) OnJoinOpen(req map[string]interface{}) (err error) {
	log.LogD(" OnJoinOpen %v", req)
	return
}

// 离开房间回调
func (d *App) OnLeaveRoom(req map[string]interface{}) (err error) {
	log.LogD(" OnLeaveRoom %v", req)
	return
}

// 踢人回调
func (d *App) OnKickPlayer(req map[string]interface{}) (err error) {
	log.LogD(" OnKickPlayer %v", req)
	return
}

// 连接状态回调
func (d *App) OnUserState(req map[string]interface{}) (err error) {
	log.LogD(" OnUserState %v", req)
	return
}

// 获取房间信息回调
func (d *App) OnRoomDetail(req *defines.MsRoomDetail) (err error) {
	log.LogD("OnRoomDetail %v", req)
	for _, v := range req.PlayersList {
		log.LogD("OnRoomDetail PlayersList %v", v)
	}
	log.LogD("OnRoomDetail WatchRoom %v", req.WatchRoom)
	return
}

// 设置房间属性回调
func (d *App) OnSetRoomProperty(req map[string]interface{}) (err error) {
	log.LogD(" OnSetRoomProperty %v", req)
	return
}

// 房间连接回调
func (d *App) OnHotelConnect(req map[string]interface{}) (err error) {
	log.LogD(" OnHotelConnect %v", req)
	return
}

// 消息广播
func (d *App) OnReceiveEvent(req *defines.MsOnReciveEvent) (err error) {
	// log.LogD(" OnReceiveEvent %v", string(req.CpProto))
	d.Example_Push(req)
	return
}

// 房间断开
func (d *App) OnDeleteRoom(req map[string]interface{}) (err error) {
	log.LogD(" OnDeleteRoom %v", req)
	return
}

// 连接房间检测回调
func (d *App) OnHotelCheckin(req map[string]interface{}) (err error) {
	log.LogD(" OnHotelCheckin %v", req)
	return
}

// 设置帧同步
func (d *App) OnSetFrameSyncRate(req *defines.MsFrameSyncRateNotify) (err error) {
	log.LogD(" OnHotelSetFrameSyncRate %v", req)
	return
}

// 帧数据更新
func (d *App) OnFrameUpdate(req *defines.FrameDataList) (err error) {
	// log.LogD(" OnFrameUpdate %v", req)
	for _, v := range req.Items {
		log.LogD(" OnFrameUpdate roomID 【%d】 length【%d】 CpProto [%s], SrcUserID [%d] , Timestamp [%d]", req.RoomID, len(req.Items), v.CpProto, v.SrcUserID, v.Timestamp)
	}
	return
}

func (d *App) Example_Push(req *defines.MsOnReciveEvent) {
	var optMap map[string]interface{}
	if err := json.Unmarshal(req.CpProto, &optMap); err != nil {
		log.LogE("event message Unmarshal error %v", err)
		return
	}

	cmd := optMap["cmd"].(string)
	log.LogD("event message [%v]", optMap)
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
		rate := optMap["frameRate"].(float64)
		enableGS := optMap["enableGS"].(float64)
		d.push.SetFrameSyncRate(req.GameID, uint32(rate), uint32(enableGS), req.RoomID)
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
	default:
	}
}

// 创建房间 示例
func (d *App) example_createRoom(gameID uint32) {
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
