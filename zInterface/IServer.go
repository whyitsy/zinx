package zInterface

// 基础 server 接口有三个基础方法：Start、Stop、Serve
type IServer interface {
	Start() // 启动服务器
	Stop()  // 停止服务器
	Serve() // 运行服务器
}
