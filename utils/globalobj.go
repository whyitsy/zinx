package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"zinx/zInterface"
)

/*
这个文件存储Zinx框架中的一切全局参数. 以各种结构体对象的形式对外提供, 在init()函数中进行初始化.
*/

type GlobalObj struct {
	// Server
	TcpServer zInterface.IServer
	IPVersion string
	Host      string
	TcpPort   int
	Name      string

	// Zinx
	Version        string
	MaxConn        int    // 当前框架允许的最大连接数
	MaxPackageSize uint32 // 当前框架允许每次传输的最大数据包大小
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	file, err := os.OpenFile("conf/zinx.json", os.O_RDONLY, 0)
	if err != nil {
		panic(err) // 要求求必须有这个配置文件
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("close file error: ", err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("read file error: ", err)
		return
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err) // 配置了参数, 格式必须正确. 可以不配置使用默认参数.
	}
}

func init() {
	// 默认参数
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		IPVersion:      "tcp4",
		Version:        "V0.4",
		Host:           "0.0.0.0",
		TcpPort:        8999,
		MaxConn:        1000,
		MaxPackageSize: 512,
	}

	// 尝试使用配置文件来覆盖默认参数
	GlobalObject.Reload()
}
