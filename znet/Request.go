package znet

import (
	"zinx/zInterface"
)

type Request struct {
	// 已经和客户端建立好的连接
	conn zInterface.IConnection
	// 客户端请求的数据
	data []byte
}

func (c *Request) GetConnection() zInterface.IConnection {
	return c.conn
}

func (c *Request) GetData() []byte {
	return c.data
}
