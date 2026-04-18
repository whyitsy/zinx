package Zinx_Test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
	"zinx/znet"
)

// go中 单元测试的文件必须以_test.go结尾.

// TestDataPack 只负责测试拆包和封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟服务器端和客户端之间的通信, 使用两个 goroutine 来模拟服务器端和客户端.
	*/

	// 模拟服务端
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		t.Fatal("server listen err:", err)
	}
	go func() {
		// 从客户端读取数据, 进行拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("accept err:", err)
			}
			go func(conn net.Conn) {
				// 先读取数据头
				dp := znet.NewDataPack()
				for {
					// 第一次从conn中读取 headData
					headData := make([]byte, dp.GetHeadLen())
					//_, err := conn.Read(headData) // Attention: Read是当前有多少就读多少, 可能出现不足GetHeadLen()的情况, 因为TCP是面向流的协议, 这种读取方式可能一次读取不到完整的数据.
					_, err := io.ReadFull(conn, headData) // ReadFull 会一直读取直到读满指定长度的数据, 就不会出现读不完整的问题了
					if err != nil {
						fmt.Println("read head err:", err)
						break
					}
					msg, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("unpack head err:", err)
						break
					}
					if msg.GetDataLen() > 0 {
						// 说明消息体有数据, 需要继续读取
						// 第二次从conn中读取数据
						data := make([]byte, msg.GetDataLen())
						//_, err := conn.Read(data)
						_, err := io.ReadFull(conn, data) // 同样的问题, 也是处理粘包问题的重要实现方案.
						if err != nil {
							fmt.Println("read data err:", err)
							break
						}
						msg.SetData(data)
						fmt.Printf("==> Recv Msg: ID=%d, len=%d, data=%s\n", msg.GetMsgId(), msg.GetDataLen(), msg.GetData())
					}
				}
				// 再根据数据头中的长度读取数据内容
			}(conn)
		}
	}()

	// 模拟客户端
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:9999")
		if err != nil {
			fmt.Println("dial err:", err)
			return
		}
		dp := znet.NewDataPack()

		// 使用两个包来模拟粘包问题

		msg1 := &znet.Message{
			Id:      1,
			DataLen: 4,
			Data:    []byte("zinx"),
		}
		data1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("pack msg1 err:", err)
			return
		}
		msg2 := &znet.Message{
			Id:      2,
			DataLen: 10,
			Data:    []byte("hello zinx"),
		}
		data2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("pack msg2 err:", err)
			return
		}

		data := append(data1, data2...) // 将两个包的数据合并成一个包, 模拟粘包问题
		_, err = conn.Write(data)
		if err != nil {
			return
		}
	}()

	// 让测试函数等待一段时间, 以便模拟的服务器和客户端能够完成通信
	timeContext := time.NewTimer(time.Second * 5)
	<-timeContext.C
}
