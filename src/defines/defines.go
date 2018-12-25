package defines

// 创建房间接口参数
type MsOnCreateRoom struct {
	GameID       uint32
	UserID       uint32
	UserProfile  []byte
	RoomID       uint64
	RoomProperty []byte
	MaxPlayer    uint32
	State        uint32
	Mode         int32
	CanWatch     int32
	CreateFlag   uint32
	CreateTime   uint64
}

// 有人加入房间接口参数
type MsOnJoinRoom struct {
	RoomID      uint64
	UserID      uint32
	GameID      uint32
	UserProfile []byte
	JoinType    uint32
	MaxPlayers  uint32
	Checkins    []uint32
	Players     []uint32
}

type MsPlayerInfo struct {
	UserID      uint32 `json:"user_id"`
	UserProfile []byte `json:"user_profile"`
}

// 观战设置类型
type MsWatchSeting struct {
	MaxWatch        uint32
	WatchPersistent bool
	WatchDelayMs    uint32
	CacheTime       uint32
}

// 游戏房间信息
type MsRoomInfo struct {
	RoomName     string `json:"roomName,omitempty"`
	MaxPlayer    uint32 `json:"maxPlayer,omitempty"`
	Mode         int32  `json:"mode,omitempty"`
	CanWatch     int32  `json:"canWatch,omitempty"`
	Visibility   int32  `json:"visibility,omitempty"`
	RoomProperty string `json:"roomProperty,omitempty"`
}

type MsCreateRoomReq struct {
	GameID   uint32         `json:"game_id"`
	Ttl      uint32         `json:"ttl"`
	RoomInfo *MsRoomInfo    `json:"room_info"`
	WatchSet *MsWatchSeting `json:"watch_set"`
}

type MsCreateRoomRsp struct {
	Status uint32 `json:"status"`
	RoomID uint64 `json:"room_id"`
}

// 观战房间信息
type MsWatchRoom struct {
	RoomID           uint64          `json:"room_id"`
	State            uint32          `json:"state"`
	CurWatch         uint32          `json:"cur_watch"`
	WatchSet         *MsWatchSeting  `json:"watchset"`
	WatchPlayersList []*MsPlayerInfo `json:"watch_player_list"`
}

type MsOnReciveEvent struct {
	GameID    uint32
	RoomID    uint64
	UserID    uint32
	Flag      uint32
	DestsList []uint32
	CpProto   []byte
}

type MsPushEventReq struct {
	PushType  int32    // 推送类型，配合destsList使用 1：推送给列表中的指定用户，2：推送给除列表中指定用户外的其他用户，3：推送给房间内的所有用户
	GameID    uint32   //游戏ID
	RoomID    uint64   //房间ID
	DestsList []uint32 //userID列表
	CpProto   []byte   //消息内容，不超过1024字节
}

type MsTeamInfoItem struct {
	TeamID     uint64
	Password   string
	Capacity   uint32
	Mode       int32
	Owner      uint32
	PlayerList []*MsPlayerInfo
}

type MsBrigadeItem struct {
	BrigadeID uint32
	TeamList  []*MsTeamInfoItem
}

// 获取房间详细信息返回参数
type MsRoomDetail struct {
	GameID       uint32           `json:"game_id"`
	RoomID       uint64           `json:"room_id"`
	UserID       uint32           `json:"user_id"`
	State        uint32           `json:"state"`
	MaxPlayer    uint32           `json:"max_player"`
	Mode         int32            `json:"mode"`
	CanWatch     int32            `json:"can_watch"`
	RoomProperty string           `json:"room_property"`
	Owner        uint32           `json:"owner"`
	CreateFlag   uint32           `json:"create_flag"`
	WatchRoom    *MsWatchRoom     `json:"watch_room"`
	PlayersList  []*MsPlayerInfo  `json:"player_list"`
	BrigadeList  []*MsBrigadeItem `json:"brigade_list"`
}

// 设置帧同步请求
type MsSetFrameSyncRateReq struct {
	RoomID       uint64
	GameID       uint32
	FrameRate    uint32
	EnableGS     uint32
	CacheFrameMS int32
}

//
type MsFrameSyncRateNotify struct {
	GameID       uint32
	RoomID       uint64
	FrameRate    uint32
	FrameIdx     uint32
	Timestamp    uint64
	EnableGS     uint32
	CacheFrameMS int32
}

// 每帧数据中包含多个数据项
type MsFrameDataItem struct {
	SrcUserID uint32
	CpProto   []byte
	Timestamp uint64
}

// 每帧数据
type MsFrameDataList struct {
	GameID     uint32
	RoomID     uint64
	FrameIndex uint32
	Items      []*MsFrameDataItem
}
