package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	// ip与端口
	Ip   string
	Port int

	// 添加字段
	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// 创建一个Server的构造方法
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 业务方法
func (s *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	// fmt.Println("连接建立成功")

	user := NewUser(conn)
	// 用户上线，将用户加入的OnlineMap中
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	// 广播当前用户上线消息
	s.BroadCast(user, "已上线")

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, "下线")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err: ", err)
				return
			}

			// 提取用户信息（去除“\n”）
			msg := string(buf[:n-1])

			// 将得到的信息进行广播
			s.BroadCast(user, msg)
		}
	}()

}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		// 将message发送给全部在线User
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	s.Message <- sendMsg
}

// 启动服务器的接口
func (s *Server) Start() {
	// 监听地址
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	// 启动监听message的goroutine
	go s.ListenMessage()

	// 阻塞等待并接受客户端的连接请求
	// 成功时返回 net.Conn 接口类型
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listen accept err: ", err)
			continue
		}

		// 若连接成功，就开一个协程去执行当前连接的业务
		go s.Handler(conn)
	}
}
