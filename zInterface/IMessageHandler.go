package zInterface

type IMessageHandler interface {
	// 进行消息 Id 到 router的映射
	DoMessageHandler(request IRequest)
	// 为消息添加路由
	AddRouter(msgId uint32, router IRouter)
	// 对外提供的api, 用于启动Worker工作池
	StartWorkerPool()
	// 将消息交给TaskQueue, 按照平均分配算法.
	SendMsgToTaskQueue(request IRequest)
}
