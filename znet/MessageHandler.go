package znet

import (
	"fmt"
	"strconv"
	"zinx/zInterface"
)

type MessageHandler struct {
	Apis map[uint32]zInterface.IRouter
}

func NewMessageHandler() zInterface.IMessageHandler {
	return &MessageHandler{
		Apis: make(map[uint32]zInterface.IRouter),
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
