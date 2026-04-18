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

	ExitChan chan bool // 由 Reader 来 通知 Writer 退出的管道, Reader知道连接是否已经关闭.

	messageChan chan []byte // 无缓冲的管道, 实现将Reader goroutine的数据传给 Writer goroutine, 由 Writer 来负责发送给客户端

	MessageHandler zInterface.IMessageHandler
}

func NewConnection(conn *net.TCPConn, connID uint32, messageHandler zInterface.IMessageHandler) *Connection {
	return &Connection{
		conn:           conn,
		connID:         connID,
		isClosed:       false,
		ExitChan:       make(chan bool, 1),
		messageChan:    make(chan []byte),
		MessageHandler: messageHandler,
	}
}

// 连接的读业务方法
func (c *Connection) startReader() {
	fmt.Println("[Reader Goroutine is running] ConnID =", c.connID)
	defer fmt.Println("[Reader exited!] ConnID =", c.connID)
	defer c.Stop() // reader 退出了, 就调用 Stop 来关闭连接, 这样就会通知 writer 退出了

	for {
		// region 拆包流程
		dp := NewDataPack()
		headerBuf := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.conn, headerBuf); err != nil {
			fmt.Println("read head error: ", err) // wsarecv 错误是Windows系统特有的错误, 在linux上时 io.EOF 错误.
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

// 将消息发送给写业务, 如果Reader 中需要处理写业务, 就将消息发送给 Writer. 这时mesID 就可以先来区分哪些消息需要处理写业务, 哪些消息只需要处理读业务了
func (c *Connection) startWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[Writer exited!] ConnID =", c.connID)

	// 等待 reader 发送消息或者等待 conn 关闭
	for {
		select {
		case data := <-c.messageChan:
			// 有数据要写给客户端, 这里是封包之后的数据了
			if _, err := c.conn.Write(data); err != nil {
				fmt.Println("writer send data error: ", err)
				return // 结束本次 writer goroutine
			}
		case <-c.ExitChan:
			// reader 已经退出了, writer 也要退出
			return
		}
	}

}

func (c *Connection) Start() {
	fmt.Println("Connection Start()... ConnID = ", c.connID)
	// 先启动从当前连接中读取数据的业务
	go c.startReader()

	go c.startWriter()

}

func (c *Connection) Stop() {
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

	// 通知 writer 退出
	c.ExitChan <- true

	// 关闭资源
	close(c.ExitChan)
	close(c.messageChan)
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

	// 将 msg 发送writer, 由 writer 来负责发送给客户端
	c.messageChan <- binaryMsg

	return nil
}
