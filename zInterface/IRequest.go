package zInterface

/*
IRequest接口：将客户端请求的连接和数据包装到一个Request中.
*/
type IRequest interface {
	// 获取当前的连接
	GetConnection() IConnection
	// 获取当前连接的数据
	GetData() []byte
}
