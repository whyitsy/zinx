package zInterface

import "net"

// 连接模块抽象层

type IConnection interface {
	// 启动, 完成连接相关的准备工作
	Start()
	// 停止连接
	Stop()
	// 获取当前连接处理的 connection 对象
	GetTCPConnection() *net.TCPConn
	// 每一个模块都有一个连接ID 以区分不同的连接
	GetConnID() uint32
	// 获取远程客户端的 TCP 状态 IP Port
	RemoteAddr() net.Addr
	// 发送数据, 将数据发送给远程的客户端
	Send(data []byte) error
}

// 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
