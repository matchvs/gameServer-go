package test

import (
	pb "commonlibs/proto"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/matchvs/gameServer-go/src/message"
)

func Test_FrameRateSet(t *testing.T) {
	var (
		hd_test       = &GsDefaultHandler{}
		cache_test    = message.NewMessageCache()
		hotelMsg_test = message.NewHotelModel(hd_test, cache_test)
	)
	rateSet := &pb.GSSetFrameSyncRate{
		GameID:    200773,
		RoomID:    14567894561454,
		Priority:  0,
		FrameRate: 10,
		FrameIdx:  1,
		EnableGS:  1,
	}
	buf, _ := proto.Marshal(rateSet)
	req := getMessagePackage(uint32(pb.HotelGsCmdID_GSSetFrameSyncRateNotifyCMDID), buf)
	hotelMsg_test.Route(200773, req)
	Route_FrameNotify(hotelMsg_test)
	Route_FrameUpdate(hotelMsg_test)
	Route_RecvEvent(hotelMsg_test)
	time.Sleep(time.Second * 20)
}

func Route_FrameNotify(ht *message.HotelMessage) {

	// aw := new(sync.WaitGroup)
	go func() {
		for i := 0; i < 500; i++ {
			go func(index uint32) {
				frame := &pb.GSFrameDataNotify{
					GameID:    200773,
					RoomID:    14567894561454,
					SrcUid:    123456,
					Priority:  0,
					CpProto:   []byte("test gameServer"),
					TimeStamp: uint64(time.Now().Unix()),
					FrameIdx:  index,
				}
				buf, _ := proto.Marshal(frame)
				req := getMessagePackage(1610, buf)

				ht.Route(200773, req)
			}(uint32(i))
			time.Sleep(time.Millisecond * 50)
		}
	}()

}

func Route_FrameUpdate(ht *message.HotelMessage) {
	go func() {
		for i := 0; i < 500; i++ {
			go func(index uint32) {
				frame := &pb.GSFrameSyncNotify{
					GameID:    200773,
					RoomID:    14567894561454,
					LastIdx:   index,
					NextIdx:   index + 1,
					Priority:  0,
					StartTS:   uint64(time.Now().Unix()),
					EndTS:     uint64(time.Now().Unix()),
					TimeStamp: uint64(time.Now().Unix()),
				}
				buf, _ := proto.Marshal(frame)
				req := getMessagePackage(uint32(pb.HotelGsCmdID_GSFrameSyncNotifyCMDID), buf)
				ht.Route(200773, req)
			}(uint32(i))
			time.Sleep(time.Millisecond * 100)
		}
	}()

}

func Route_RecvEvent(ht *message.HotelMessage) {
	go func() {
		for i := 0; i < 5000; i++ {
			go func(index uint32) {
				data := &struct {
					Cmd  string `json:"cmd"`
					Data string `json:"data"`
				}{Cmd: "event", Data: "gameServer test uint"}
				// fmt.Printf("data:%v", *data)
				databuf, _ := json.Marshal(data)
				for j := 0; j <= 20; j++ {
					if j >= 20 {
						data.Cmd = "end"
						databuf, _ = json.Marshal(data)
					}
					event := &pb.HotelBroadcast{
						GameID:  200773,
						RoomID:  14567894561454,
						Flag:    1,
						DstUids: []uint32{123456},
						CpProto: databuf,
					}
					buf, _ := proto.Marshal(event)
					req := getMessagePackage(uint32(pb.HotelGsCmdID_HotelBroadcastCMDID), buf)
					ht.Route(200773, req)
					time.Sleep(time.Millisecond * 100)
				}
				return
			}(uint32(i))
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

// get a test package which simulate server data package
func getMessagePackage(cmdID uint32, msg []byte) *pb.Package_Frame {
	pkg := new(pb.Package_Frame)
	pkg.CmdId = cmdID
	pkg.Version = 2
	pkg.Reserved = 234
	pkg.Type = pb.Package_LeagueMessage
	pkg.Message = msg[:]
	return pkg
}
