package zInterface

type IMessage interface {
	// 获取消息 ID
	GetMsgId() uint32
	// 获取消息长度
	GetDataLen() uint32
	// 获取消息内容
	GetData() []byte
	// 设置消息 ID
	SetMsgId(uint32)
	// 设置消息长度
	SetDataLen(uint32)
	// 设置消息内容
	SetData([]byte)
}
