package message

import (
	"github.com/matchvs/gameServer-go/src/defines"
)

type IHandler interface {
	// 创建房间回调
	OnCreateRoom(map[string]interface{}) error
	// 加入房间回调
	OnJoinRoom(map[string]interface{}) error
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
