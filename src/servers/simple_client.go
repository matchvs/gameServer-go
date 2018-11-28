package servers

import (
	"fmt"

	"github.com/matchvs/gameServer-go/src/log"
	pb "github.com/matchvs/gameServer-go/src/proto"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// SimpleClient grpc客户端
type SimpleClient struct {
	serverAddr string
	conn       *grpc.ClientConn
	client     pb.SimpleServiceClient
}

// NewSimpleClient 构造函数
func NewSimpleClient(addr string) (client *SimpleClient) {
	client = &SimpleClient{}
	client.serverAddr = addr
	// connect
	// gen new client

	return
}

func (c *SimpleClient) Run() (err error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(c.serverAddr, opts...)
	if err != nil {
		log.LogE("grpc.Dial failed:", err)
		return
	}
	c.conn = conn
	c.client = pb.NewSimpleServiceClient(conn)
	return
}

// Close 断开连接
func (c *SimpleClient) Close() {
	c.conn.Close()
}

// Send 发送消息
func (c *SimpleClient) Send(req *pb.Package_Frame) (rep *pb.Package_Frame, err error) {
	rep, err = c.client.SimpleRequest(context.Background(), req)
	if err != nil {
		fmt.Println("SimpleRequest failed")
	}
	return
}
