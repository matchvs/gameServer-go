/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-18 18:18:45
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-20 15:00:26
 * @Description: message push
 */
package message

import (
	pb "commonlibs/proto"
	"commonlibs/servers"

	"github.com/golang/protobuf/proto"
	"github.com/matchvs/gameServer-go/src/defines"
)

// PushManager 消息推送管理类型
type PushManager struct {
	adaptor *GSAdaptor
	roomMg  *servers.RoomManager
}

func NewPushManager(adaptor *GSAdaptor, roomMg *servers.RoomManager) (p *PushManager) {
	p = new(PushManager)
	p.adaptor = adaptor
	p.roomMg = roomMg
	return
}

func (self *PushManager) PushEvent(req *defines.MsPushEventReq) error {
	event := &pb.PushToHotelMsg{
		PushType: pb.PushMsgType(req.PushType),
		GameID:   req.GameID,
		RoomID:   req.RoomID,
		DstUids:  req.DestsList[:],
		CpProto:  req.CpProto[:],
	}
	msg, _ := proto.Marshal(event)
	return self.adaptor.PushHotel(uint32(pb.HotelGsCmdID_HotelPushCMDID), req.RoomID, msg)
}

func (self *PushManager) JoinOver(gameID uint32, roomID uint64) {
	req := new(pb.JoinOverReq)
	req.GameID = gameID
	req.RoomID = roomID
	msg, _ := proto.Marshal(req)
	self.adaptor.PushMvs(uint32(pb.MvsGsCmdID_MvsJoinOverReq), msg)
}

func (self *PushManager) JoinOpen(gameID uint32, roomID uint64) {
	req := new(pb.JoinOpenReq)
	req.GameID = gameID
	req.RoomID = roomID
	msg, _ := proto.Marshal(req)
	self.adaptor.PushMvs(uint32(pb.MvsGsCmdID_MvsJoinOpenReq), msg)
}

func (self *PushManager) KickPlayer(destID uint32, roomID uint64) {
	req := new(pb.KickPlayer)
	req.RoomID = roomID
	req.UserID = destID
	msg, _ := proto.Marshal(req)
	self.adaptor.PushMvs(uint32(pb.MvsGsCmdID_MvsKickPlayerReq), msg)
}

func (self *PushManager) GetRoomDetail(gameID, latestWatcherNum uint32, roomID uint64) {
	req := new(pb.GetRoomDetailReq)
	req.GameID = gameID
	req.RoomID = roomID
	req.LatestWatcherNum = latestWatcherNum
	msg, _ := proto.Marshal(req)
	self.adaptor.PushMvs(uint32(pb.MvsGsCmdID_MvsGetRoomDetailReq), msg)
}

func (self *PushManager) SetRoomProperty(gameID uint32, roomID uint64, roomProperty string) {
	req := new(pb.SetRoomPropertyReq)
	req.GameID = gameID
	req.RoomID = roomID
	req.RoomProperty = []byte(roomProperty)
	msg, _ := proto.Marshal(req)
	self.adaptor.PushMvs(uint32(pb.MvsGsCmdID_MvsSetRoomPropertyReq), msg)
}

// CreateRoom 主动创建房间
// crtm ： 创建房间的参数信息
// 返回类型 MsCreateRoomRsp 是创建后的状态信息
func (self *PushManager) CreateRoom(crtm *defines.MsCreateRoomReq) (*defines.MsCreateRoomRsp, error) {
	req := &pb.CreateRoom{}

	req.RoomInfo = &pb.RoomInfo{
		RoomName:     crtm.RoomInfo.RoomName,
		RoomProperty: []byte(crtm.RoomInfo.RoomProperty),
		CanWatch:     crtm.RoomInfo.CanWatch,
		MaxPlayer:    crtm.RoomInfo.MaxPlayer,
		Mode:         crtm.RoomInfo.Mode,
		Visibility:   crtm.RoomInfo.Visibility,
	}
	req.WatchSetting = &pb.WatchSetting{
		MaxWatch:        crtm.WatchSet.MaxWatch,
		WatchPersistent: crtm.WatchSet.WatchPersistent,
		WatchDelayMs:    crtm.WatchSet.WatchDelayMs,
		CacheTime:       crtm.WatchSet.CacheTime,
	}
	req.GameID = crtm.GameID
	req.Ttl = crtm.Ttl

	ack, err := self.roomMg.CreateRoom(req)
	if err != nil {
		return nil, err
	}
	rsp := &defines.MsCreateRoomRsp{
		Status: ack.Status,
		RoomID: ack.RoomID,
	}
	return rsp, nil
}

// TouchRoom 设置房间的存活时间
// 返回 200 表示成功
// gameID : 游戏ID
// ttl 	  : 空房间存活时长(房间没有任何玩家的情况)，单位秒，最大取值 86400 秒（1天）
// roomID : 房间ID
func (self *PushManager) TouchRoom(gameID, ttl uint32, roomID uint64) (uint32, error) {
	req := new(pb.TouchRoom)
	req.GameID = gameID
	req.RoomID = roomID
	req.Ttl = ttl
	ack, err := self.roomMg.TouchRoom(req)
	if err != nil {
		return 0, err
	}
	return ack.Status, nil
}

// DestroyRoom 主动销毁房间，可以销毁任意房间，如果房间中有人，就会把人剔出房间
// 返回 200 表示成功
// gameID : 游戏ID
// roomID ：房间ID
func (self *PushManager) DestroyRoom(gameID uint32, roomID uint64) (uint32, error) {
	req := new(pb.DestroyRoom)
	req.RoomID = roomID
	req.GameID = gameID
	ack, err := self.roomMg.DestroyRoom(req)
	if err != nil {
		return 0, err
	}
	return ack.Status, nil
}

// SetFrameSyncRate 设置帧同步参数
// gameID : 游戏ID
// frameRate : 帧率（0到20，且能被1000整除）
// enableGS GameServer是否参与帧同步（0：不参与；1：参与）
// roomID : 要设置帧同步的房间ID
func (self *PushManager) SetFrameSyncRate(setinfo *defines.MsSetFrameSyncRateReq) {
	req := new(pb.GSSetFrameSyncRate)
	req.RoomID = setinfo.RoomID
	req.GameID = setinfo.GameID
	req.FrameRate = setinfo.FrameRate
	req.FrameIdx = 1
	req.Priority = 0
	req.EnableGS = setinfo.EnableGS
	req.CacheFrameMS = setinfo.CacheFrameMS
	msg, _ := proto.Marshal(req)
	self.adaptor.PushHotel(uint32(pb.HotelGsCmdID_GSSetFrameSyncRateCMDID), req.RoomID, msg)
}

// FrameBroadcast 发送帧同步数据给 游戏房间服务
// gameID : 游戏ID
// operation : 数据处理方式 0：只发客户端；1：只发GS；2：同时发送客户端和GS
// roomID : 房间ID
// cpProto : 要发送的数据
func (self *PushManager) FrameBroadcast(gameID uint32, operation int32, roomID uint64, cpProto []byte) {
	req := new(pb.GSFrameBroadcast)
	req.Priority = 0
	req.GameID = gameID
	req.RoomID = roomID
	req.Operation = operation
	req.CpProto = cpProto[:len(cpProto)]
	msg, _ := proto.Marshal(req)
	self.adaptor.PushHotel(uint32(pb.HotelGsCmdID_GSFrameBroadcastCMDID), req.RoomID, msg)
}

// GetOffLineCacheData 获取断线重新连接后的缓存数据
// gameID : 游戏ID
// cacheMS : 想要获取的毫秒数(-1表示获取所有缓存数据，该字段的赋值上限为1小时)
// roomID : 房间ID
func (self *PushManager) GetOffLineCacheData(gameID uint32, roomID uint64, cacheMS int32) error {
	req := new(pb.GSGetCacheData)
	req.GameID = gameID
	req.RoomID = roomID
	req.CacheFrameMS = cacheMS
	msg, err := proto.Marshal(req)
	if err != nil {
		return err
	}
	return self.adaptor.PushHotel(uint32(pb.HotelGsCmdID_GSGetCacheDataCMDID), req.RoomID, msg)
}
