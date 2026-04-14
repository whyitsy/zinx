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

	handleAPI zInterface.HandleFunc

	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callback zInterface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callback,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
}

// 连接的读业务方法
func (c *Connection) startReader() {
	fmt.Println("start reader goroutine... ConnID = ", c.ConnID)
	defer fmt.Println("connID = ", c.ConnID, " reader is exit, remote addr is ", c.RemoteAddr())
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error: ", err)
			continue
		}
		// 调用当前连接所绑定的业务方法, 处理客户端请求的消息
		if err := c.handleAPI(c.Conn, buf[:cnt], cnt); err != nil {
			fmt.Println("connID = ", c.ConnID, " handle is error: ", err)
			break
		}
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
