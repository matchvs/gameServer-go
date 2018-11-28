/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-11-27 20:08:05
 * @LastEditors: Ville
 * @LastEditTime: 2018-11-28 14:35:56
 * @Description: file content
 */

package servers

import (
	pb "commonlibs/proto"
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/matchvs/gameServer-go/src/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	bufferLen         = 100000
	reconnectTimeout  = 10
	maxReserved       = 1000
	sendTimeout       = 10
	heartBeatInterval = 5
)

type PushHandler interface {
	RecvPush(*pb.Package_Frame) error
}

type StreamClient struct {
	serverAddr  string
	serverID    uint32
	load        uint32
	timeout     int64
	stop        int32
	isconnect   bool
	reserved    uint64
	sendCh      chan *pb.Package_Frame
	recvCh      chan int
	conn        *grpc.ClientConn
	stream      pb.CSStream_StreamClient
	rwlock      sync.RWMutex
	pushHandler PushHandler
}

func NewStreamClient(addr string, serverID uint32, ph PushHandler) (sc *StreamClient) {
	sc = new(StreamClient)
	sc.serverAddr = addr
	sc.serverID = serverID
	sc.stop = 0
	sc.isconnect = false
	sc.reserved = 1
	sc.sendCh = make(chan *pb.Package_Frame, bufferLen)
	sc.pushHandler = ph
	sc.timeout = 10
	if serverID == 0 {
		sc.serverID = crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s", addr)))
	}
	return
}

func (s *StreamClient) Start() error {
	err := s.connect()
	if err != nil {
		log.LogE("stream client start connect remote addr %s error %v", s.serverAddr, err)
		return err
	}
	s.onConnect()

	var wg = new(sync.WaitGroup)
	wg.Add(2)
	go s.sendRoutine(wg)
	go s.recvRoutine(wg)
	wg.Wait()
	return nil
}

func (s *StreamClient) Stop() {
	atomic.AddInt32(&s.stop, 1)
	s.disconnect()
	close(s.sendCh)
	close(s.recvCh)
}

func (s *StreamClient) IsStop() int32 {
	stop := atomic.LoadInt32(&s.stop)
	return stop
}

func (s *StreamClient) connect() (err error) {
	var (
		opts []grpc.DialOption
	)
	opts = append(opts, grpc.WithInsecure())
	s.conn, err = grpc.Dial(s.serverAddr, opts...)
	if err != nil {
		log.LogE("dail error: %v", err)
		return
	}

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("userid", fmt.Sprintf("%d", s.serverID)),
	)

	client := pb.NewCSStreamClient(s.conn)
	s.stream, err = client.Stream(ctx)
	if err != nil {
		s.conn.Close()
		s.conn = nil
		log.LogE("connect server error: %v", err)
		return err
	}
	return
}

func (s *StreamClient) onConnect() {
	s.rwlock.Lock()
	s.isconnect = true
	s.rwlock.Unlock()
}

func (s *StreamClient) ConnectSuccess() (ok bool) {
	s.rwlock.Lock()
	ok = s.isconnect
	s.rwlock.Unlock()
	return ok
}

// 断开连接
func (s *StreamClient) disconnect() {
	s.rwlock.Lock()
	s.isconnect = false
	s.rwlock.Unlock()
}

func (s *StreamClient) GetReserved() uint32 {
	s.reserved++
	if s.reserved > maxReserved {
		s.reserved = 0
	}
	return uint32(s.reserved)
}

func (s *StreamClient) Push(cmdid uint32, msg []byte) error {
	resp := &pb.Package_Frame{
		Type:     pb.Package_LeagueMessage,
		Version:  VERSION,
		CmdId:    cmdid,
		UserId:   uint64(s.serverID),
		Reserved: s.GetReserved(),
		Message:  msg,
	}
	s.send(resp)
	return nil
}

// 发送消息
func (s *StreamClient) send(frame *pb.Package_Frame) (err error) {
	if s.ConnectSuccess() == false {
		log.LogW("no connect server (%v) ", s.serverAddr)
		err = errors.New("no connect")
		return
	}
	select {
	case s.sendCh <- frame:
		break
	case <-time.After(time.Second * time.Duration(s.timeout)):
		log.LogE("send frame to server (%v) failed, timeout", s.serverAddr)
		s.disconnect()
		return errors.New("send message timeout")
	}
	return
}

func (s *StreamClient) connRoutine(wg *sync.WaitGroup) {
}

// 轮询是否有消息需要发送
func (s *StreamClient) sendRoutine(wg *sync.WaitGroup) {
	wg.Done()
	for {
		frame, ok := <-s.sendCh
		if ok {
			err := s.stream.Send(frame)
			if err != nil {
				log.LogE("send frame to server [%v] error: %v", s.serverAddr, err)
			}
		} else {
			log.LogD("sendRoutine close, server [%v]", s.serverAddr)
			s.disconnect()
			return
		}
	}
}

// 轮询是否接收消息
func (s *StreamClient) recvRoutine(wg *sync.WaitGroup) {
	wg.Done()
	for {
		select {
		case <-s.recvCh:
			log.LogD("recv a close operation toServerAddr [%v]", s.serverAddr)
			return
		default:
			break
		}

		frame, err := s.stream.Recv()
		if err == io.EOF {
			log.LogW("server close stream : %v", err)
			s.disconnect()
		}
		if err != nil {
			log.LogE("stream client recive message error %v", err)
			s.disconnect()
			return
		}
		if s.pushHandler != nil && frame != nil {
			s.pushHandler.RecvPush(frame)
		} else {
			log.LogD("stream client recive message success [%v]", frame.CmdId)
		}
	}
}
