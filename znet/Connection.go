package znet

import (
	"fmt"
	"net"
	"zinx/zInterface"
)

type Connection struct {
	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	ExitChan chan bool // 告知当前连接已经退出/停止的 channel

	Router zInterface.IRouter // 该连接处理的方法Router
}

func NewConnection(conn *net.TCPConn, connID uint32, router zInterface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
}

// 连接的读业务方法
func (c *Connection) startReader() {
	fmt.Println("start reader goroutine... ConnID = ", c.ConnID)
	defer fmt.Println("connID = ", c.ConnID, " reader is exit, remote addr is ", c.RemoteAddr())
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error: ", err)
			continue
		}

		req := Request{
			conn: c,
			data: buf,
		}

		// 这里是用 goroutine 来处理请求. TODO: 这里用另一个 goroutine 来处理请求 会比继续使用当前的 goroutine 来处理更好吗？
		go func(request zInterface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

func (c *Connection) Start() {
	fmt.Println("Connection Start()... ConnID = ", c.ConnID)
	// 先启动从当前连接中读取数据的业务
	go c.startReader()
	// TODO 后面需要继续完善从当前连接写数据的业务

}

func (c *Connection) Stop() {
	fmt.Println("Connection Stop()... ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭 socket 连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("close tcp Conn error: ", err)
		return
	}

	// 关闭资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}
