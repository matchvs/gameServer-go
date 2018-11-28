package message

import (
	pb "commonlibs/proto"
	"commonlibs/servers"
	"fmt"
	"sync"
)

type MessageCache struct {
	JoinRoom map[string]*pb.Request
	lock     sync.RWMutex
}

func NewMessageCache() *MessageCache {
	cache := new(MessageCache)
	cache.JoinRoom = make(map[string]*pb.Request)
	return cache
}
func (self *MessageCache) AddWaitJoin(roomID uint64, userID uint32, val *pb.Request) {
	self.lock.Lock()
	self.JoinRoom[fmt.Sprintf("%d_%d", roomID, userID)] = val
	self.lock.Unlock()
}
func (self *MessageCache) DelWaitJoin(roomID uint64, userID uint32) {
	self.lock.Lock()
	if len(self.JoinRoom) > 0 {
		delete(self.JoinRoom, fmt.Sprintf("%d_%d", roomID, userID))
	}
	self.lock.Unlock()
}
func (self *MessageCache) GetWaitJoin(roomID uint64, userID uint32) (val *pb.Request) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if len(self.JoinRoom) > 0 {
		val = self.JoinRoom[fmt.Sprintf("%d_%d", roomID, userID)]
	}
	return
}

type MessageModel struct {
	handle   IHandler
	msgCache *MessageCache
}

func (m *MessageModel) CanDeal(cmdid int32) bool {
	return false
}

func (m *MessageModel) Route(connID uint64, req servers.GSRequest) (res []byte, err error) {
	res = []byte("hello")
	return
}
