package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/zInterface"
)

type MessageHandler struct {
	Apis map[uint32]zInterface.IRouter
	// 消息队列 channel, 由Worker负责取出消息进行处理
	TaskQueue []chan zInterface.IRequest
	// 业务工作 Worker池的 Worker数量, 应该与TaskQueue的数量一致.
	WorkerPoolSize uint32
}

func NewMessageHandler() zInterface.IMessageHandler {
	return &MessageHandler{
		Apis:           make(map[uint32]zInterface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan zInterface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MessageHandler) DoMessageHandler(request zInterface.IRequest) {
	// 1. 从Request中找到msgId
	msgId := request.GetMsgID()
	router, ok := mh.Apis[msgId]
	if !ok {
		fmt.Println("api msgId = ", msgId, " is not FOUND! Need Call AddRouter() First!") // Golang没有nameof方法，有点遗憾.
		return
	}
	// 2. 根据msgId找到对应的router并执行
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (mh *MessageHandler) AddRouter(messageId uint32, router zInterface.IRouter) {
	if _, ok := mh.Apis[messageId]; ok {
		// 当前api已经存在, 直接panic
		panic("repeated api, messageId = " + strconv.Itoa(int(messageId)))
	}
	mh.Apis[messageId] = router
	fmt.Printf("MessageId %d Router Add Succeed!\n", messageId)
}

// StartWorkerPool 对外提供的api, 调用这个方法按配置数量初始化Worker池.
func (mh *MessageHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 创建worker使用的channel，使用序号将channel与worker绑定
		mh.TaskQueue[i] = make(chan zInterface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前Worker，阻塞等待消息从channel传递过来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 这个不对外暴露, 可以是小写开头, 也可以不写到IMessageHandler接口中.
func (mh *MessageHandler) StartOneWorker(workerId int, taskQueue chan zInterface.IRequest) {
	fmt.Printf("Worker ID = %d is started.\n", workerId)
	// 不断的等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来, 就取出队列的Request, 执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMessageHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息发送到消息队列. 简单的平均分配算法.
func (mh *MessageHandler) SendMsgToTaskQueue(request zInterface.IRequest) {
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Printf("Add Request MsgID = %d, ConnID = %d to WorkerID = %d\n",
		request.GetMsgID(), workerId,
		request.GetConnection().GetConnID())
	mh.TaskQueue[workerId] <- request
}
