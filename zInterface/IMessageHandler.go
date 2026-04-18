package zInterface

type IMessageHandler interface {
	// 进行消息 Id 到 router的映射
	DoMessageHandler(request IRequest)
	// 为消息添加路由
	AddRouter(msgId uint32, router IRouter)
}
