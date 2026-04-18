package znet

import (
	"fmt"
	"io"
	"net"
	"zinx/zInterface"
)

type Connection struct {
	conn *net.TCPConn

	connID uint32

	isClosed bool

	ExitChan chan bool // 告知当前连接已经退出/停止的 channel

	MessageHandler zInterface.IMessageHandler
}

func NewConnection(conn *net.TCPConn, connID uint32, messageHandler zInterface.IMessageHandler) *Connection {
	return &Connection{
		conn:           conn,
		connID:         connID,
		isClosed:       false,
		ExitChan:       make(chan bool, 1),
		MessageHandler: messageHandler,
	}
}

// 连接的读业务方法
func (c *Connection) startReader() {
	fmt.Println("start reader goroutine... ConnID = ", c.connID)
	defer fmt.Println("connID =", c.connID, " reader is exit, remote addr is ", c.RemoteAddr())
	defer c.Stop()

	for {
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//cnt, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read buf error: ", err)
		//	break
		//}

		// region 拆包流程
		dp := NewDataPack()
		headerBuf := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.conn, headerBuf)
		if err != nil {
			fmt.Println("read head error: ", err)
			break
		}
		message, err := dp.UnPack(headerBuf)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}
		if message.GetDataLen() > 0 {
			dataBuf := make([]byte, message.GetDataLen())
			_, err := io.ReadFull(c.conn, dataBuf)
			if err != nil {
				fmt.Println("read data error: ", err)
				break
			}
			message.SetData(dataBuf)
			//fmt.Printf("==> Recv Msg: ID=%d, len=%d, data=%s\n", message.GetMsgId(), message.GetDataLen(), message.GetData())
		}
		// endregion

		req := Request{
			conn: c,
			msg:  message,
		}

		// 这里是用 goroutine 来处理请求. TODO: 这里用另一个 goroutine 来处理请求 会比继续使用当前的 goroutine 来处理更好吗？
		go c.MessageHandler.DoMessageHandler(&req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Connection Start()... ConnID = ", c.connID)
	// 先启动从当前连接中读取数据的业务
	go c.startReader()
	// TODO 后面需要继续完善从当前连接写数据的业务

}

func (c *Connection) Stop() {
	fmt.Println("Connection Stop... ConnID = ", c.connID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭 socket 连接
	err := c.conn.Close()
	if err != nil {
		fmt.Println("close tcp Conn error: ", err)
		return
	}

	// 关闭资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) GetConnID() uint32 {
	return c.connID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SendMsg 对发送的数据进行封包, 然后发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	// 判断连接是否关闭
	if c.isClosed == true {
		return fmt.Errorf("connection closed when send msg")
	}

	dp := NewDataPack()
	message := NewMessage(msgId, data)
	binaryMsg, err := dp.Pack(message)
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return err
	}

	// 将 msg 发送给客户端
	_, err = c.conn.Write(binaryMsg)
	if err != nil {
		return err
	}

	return nil
}
