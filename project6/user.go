package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	// 将客户端当前连接的地址作为地址
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// 启动当前监听 user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 创建用户上线的业务
func (u *User) Online() {
	// 用户上线，将用户加入的OnlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// 广播当前用户上线消息
	u.server.BroadCast(u, "已上线")
}

// 创建用户下线的业务
func (u *User) Offline() {
	// 用户下线，将用户从OnlineMap中删除
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// 广播当前用户上线消息
	u.server.BroadCast(u, "下线")
}

// 用户处理消息业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + "在线\n"
			// 给当前user对应的客户端发消息
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" { // 更改用户名业务，如果用户输入字符长度 > 7并且前七个字符是 rename| 则启动修改用户名业务
		// 消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]

		// 判断name是否存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("当前用户名被使用\n")
		} else {
			u.server.mapLock.Lock()
			// 更新服务端名称
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			// 更新本地用户名
			u.Name = newName
			// 向用户广播
			u.SendMsg("您已经更新用户名：" + u.Name + "\n")
		}
	} else {
		// 查询当前用户有哪些
		u.server.BroadCast(u, msg)
	}

}

// 给当前user对应的客户端发消息的方法

func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// 监听当前User channel的方法，一旦有消息，就发送给对端客户端
func (u *User) ListenMessage() {
	for msg := range u.C {
		// 从管道中取出消息，等价于
		// msg := <-u.C
		// 转换成二进制发送给对端客户端（发送后换行）
		u.conn.Write([]byte(msg + "\n"))
	}
}
