/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-11-27 20:08:05
 * @LastEditors: Ville
 * @LastEditTime: 2018-11-28 14:36:04
 * @Description: the module of server communication with other
 */

package servers

import (
	"commonlibs/errors"
	pb "commonlibs/proto"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/matchvs/gameServer-go/src/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	VERSION = 2 //通信协议版本
)

type GSRequest *pb.Package_Frame
type GSWrite func(*pb.Package_Frame) error

type IGSServer interface {
	OnConnect(userid uint64, token string, write GSWrite) error
	OnDisconnect(userid uint64, token string) error
	Route(userid uint64, req GSRequest, rsp GSWrite) (err error)
}

// 继承 pb.CSStreamServer
type StreamServer struct {
	addr       string
	timeOut    int64
	grpcServer *grpc.Server
	gsServer   IGSServer
}

func NewStreamServer(add string, gs IGSServer, timeout int64) (sc *StreamServer) {
	sc = new(StreamServer)
	sc.grpcServer = grpc.NewServer()
	sc.gsServer = gs
	sc.addr = add
	sc.timeOut = timeout
	pb.RegisterCSStreamServer(sc.grpcServer, sc)
	return
}

// 启动一个 grpc 服务端
func (s *StreamServer) Start() (err error) {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.LogE("server listen error %v", err)
		return err
	}
	log.LogI("server addr is %v", s.addr)
	s.grpcServer.Serve(lis)
	return nil
}

// 停止 grpc 服务
func (s *StreamServer) Stop() {
	s.grpcServer.Stop()
}

//获取元组数据
func getMetadata(md metadata.MD) (userid uint64, token string, err error) {
	if len(md["userid"]) == 0 {
		log.LogD("cannot read key:userid from metadata")
		return
	}
	// userid & token
	uid, err := strconv.Atoi(md["userid"][0])
	userid = uint64(uid)
	if err != nil {
		log.LogE("atoi user id failed %v", err)
		return
	}
	token = ""
	if len(md["token"]) != 0 {
		token = md["token"][0]
	}

	return
}

// 继承 proto 协议中 CSStreamServer 双向流接口
// 有客户端发送消息的时候就是走的这里
func (s *StreamServer) Stream(stream pb.CSStream_StreamServer) error {
	defer errors.PrintPanicStack()

	// recv metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.LogE("cannot read metadata from context")
		return nil
	}

	userid, token, err := getMetadata(md)
	if err != nil {
		log.LogE("cannot read metadata from context")
		return err
	}

	s.gsServer.OnConnect(userid, token, GSWrite(stream.Send))
	defer s.gsServer.OnDisconnect(userid, token)

	// create recv goroutine
	var (
		sessionCh = make(chan *pb.Package_Frame)
		stopCh    = make(chan *pb.Package_Frame)
		wg        sync.WaitGroup
	)
	wg.Add(1)
	go s.recv(stream, sessionCh, &wg)
	wg.Wait()

	for {
		select {
		case frame, ok := <-sessionCh:
			if !ok {
				close(stopCh)
				return nil
			}
			go s.route(frame, stopCh, userid, GSWrite(stream.Send))

		case <-time.After(time.Second * time.Duration(s.timeOut)):
			log.LogW("client timeout")
			close(stopCh)
			return nil
		}
	}
	// return nil
}

func (s *StreamServer) route(frame *pb.Package_Frame, stopCh chan *pb.Package_Frame, connID uint64, write GSWrite) error {
	select {
	case <-stopCh:
		return nil
	default:
	}

	if err := s.gsServer.Route(connID, frame, write); err != nil {
		return err
	}
	return nil
}

//接收客户端发送过来的消息
func (s *StreamServer) recv(stream pb.CSStream_StreamServer, sessionCh chan *pb.Package_Frame, wg *sync.WaitGroup) {
	defer func() {
		close(sessionCh)
	}()
	wg.Done()
	for {
		in, err := stream.Recv()
		// log.LogD("收到消息：%v", in)
		if err == io.EOF {
			log.LogD(" stream recive EOF")
			return
		}
		if err != nil {
			log.LogD(" stream recive error %v", err)
			return
		}
		sessionCh <- in
	}
}
