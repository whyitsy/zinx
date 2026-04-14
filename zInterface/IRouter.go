package zInterface

type IRouter interface {
	// PreHandle 在处理当前业务之前的钩子方法Hook
	PreHandle(request IRequest)

	// Handle 处理当前业务的主方法
	Handle(request IRequest)

	// PostHandle 在处理当前业务之后的钩子方法Hook
	PostHandle(request IRequest)
}
