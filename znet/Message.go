package znet

/*
将从客户端读取到的数据封装到 Message中
*/

type Message struct {
	// 消息 ID
	Id uint32
	// 消息长度
	DataLen uint32
	// 消息内容
	Data []byte // 切片
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(msgId uint32) {
	m.Id = msgId
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func NewMessage(msgId uint32, data []byte) *Message {
	return &Message{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
