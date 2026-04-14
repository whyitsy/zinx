package znet

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"zinx/zInterface"
)

// IServer 接口的实现, 实现一个Server 模块

type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的 IP 版本
	IPVersion string
	// 服务器监听的 IP
	IP string
	// 服务器监听的端口
	Port int
}

// 每个新建连接需要绑定的处理函数, 这里写死, 后面应该是开发者传入.
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Printf("CallBackToClient: receive from client data: %s, cnt = %d\n", string(data), cnt)
	if _, err := conn.Write(data); err != nil {
		fmt.Println("CallBackToClient Write error: ", err)
		return errors.New("CallBackToClient Write error")
	}
	return nil
}

func (s *Server) Start() {
	// 启动一个单体服务器需要的步骤:
	fmt.Println("[Start]-", s.Name, " Server Listening at IP: ", s.IP, " Port: ", s.Port)

	// 启动一个 goroutine 来处理服务器的业务, 这样就不会阻塞后续的 Stop 和 Serve 方法
	go func() {
		// 1. 获取一个 TCP 的 Addr
		tcpAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2. 监听服务器的地址 listen
		tcpListener, err := net.ListenTCP(s.IPVersion, tcpAddr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " error: ", err)
			return
		}
		fmt.Println("start Zinx server succeed, listening...")
		// 3. accept 客户端的连接请求, 阻塞的等待客户端连接, 处理业务(读写)
		// 不断地循环处理客户端的连接请求
		connID := uint32(0)
		for {
			conn, err := tcpListener.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp error: ", err)
				return
			}

			// 将新连接conn与callback方法进行绑定, 得到我们自己分装的连接模块
			c := NewConnection(conn, connID, CallBackToClient)
			connID++

			go c.Start()

		}
	}()

}

func (s *Server) Stop() {
	// TODO: 主要是做资源清理的工作，资源、状态、链接等等
}

func (s *Server) Serve() {
	// 不直接将 Start、Stop方法暴露给使用框架的用户, 而是使用 Serve 方法来启动服务器并阻塞, 把所有逻辑封装好.
	s.Start()

	// TODO: 这里做额外业务的扩展

	// 阻塞, 让服务器一直运行
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig // 阻塞直到接收到系统中断信号
	fmt.Println("程序正常退出, 已清理资源")
}

// NewServer 初始化 Server 模块的方法
func NewServer(name string) zInterface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
}
