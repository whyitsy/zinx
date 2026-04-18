package znet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"zinx/utils"
	"zinx/zInterface"
)

/*
	封包拆包模块
	面向TCP协议中的数据流传输, 处理TCP粘包问题
	定义每个数据的前8个字节为数据头, 用于存放数据类型和长度(两个uint32), 后续为数据内容
	TLV格式, 采用大端字节序, 即数据类型和长度的高位字节在前, 低位数据在后
*/

type DataPack struct {
}

func (dp *DataPack) GetHeadLen() uint32 {
	//数据头长度为8字节, 后4字节存放数据长度, 前4字节存放数据类型
	return 8
}

// Pack 封包方法(封装数据) 方法：Type|DataLen|Data
func (dp *DataPack) Pack(message zInterface.IMessage) ([]byte, error) {
	// Attention: 因为涉及多字节传输的问题, 需要考虑字节序问题, 即大端小端的问题. 这里采用大端字节序.
	data := bytes.NewBuffer([]byte{}) // buffer与字节数组不同, buffer是一个可变的字节数组, 可以动态添加数据, 而字节数组是固定长度的, 需要预先定义长度.
	// binary 包中操作大端小端的函数.
	err := binary.Write(data, binary.BigEndian, message.GetMsgId())
	if err != nil {
		fmt.Println("binary.Write MsgId failed:", err)
		return nil, err
	}
	err = binary.Write(data, binary.BigEndian, message.GetDataLen())
	if err != nil {
		fmt.Println("binary.Write DataLen failed:", err)
		return nil, err
	}
	err = binary.Write(data, binary.BigEndian, message.GetData())
	if err != nil {
		fmt.Println("binary.Write Data failed:", err)
		return nil, err
	}
	return data.Bytes(), nil
}

// UnPack 拆包方法: 这里传入的data 是从连接中读取的GetHeadLen() 长度的字节数组, 需要按照约定好的格式和大小端序来解析数据. 外部拿到数据头后, 就可以根据数据头中的长度信息来读取数据内容了.
func (dp *DataPack) UnPack(data []byte) (zInterface.IMessage, error) {
	// 从二进制数据中读取数据类型和长度, 然后根据长度读取数据内容
	// Attention: io.Reader和io.Writer接口的实现, 除了特殊的实现外, 一般对其进行读写操作时都会自动记录当前读写位置, 因此在读取数据时不需要手动维护读写位置, 只需要按照顺序读取数据即可.
	buffer := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(buffer, binary.BigEndian, &msg.Id); err != nil {
		fmt.Println("binary.Read MsgId failed:", err)
		return nil, err
	}
	if err := binary.Read(buffer, binary.BigEndian, &msg.DataLen); err != nil {
		fmt.Println("binary.Read DataLen failed:", err)
		return nil, err
	}

	// 判断数据长度是否超过最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		fmt.Println("Too large msg data recv, dataLen=", msg.DataLen)
		return nil, fmt.Errorf("too large msg data recv, dataLen=%d", msg.DataLen)
	}

	return msg, nil

}

func NewDataPack() zInterface.IDataPack {
	return &DataPack{}
}
