package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	// 将客户端当前连接的地址作为地址
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动当前监听 user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前User channel的方法，一旦有消息，就发送给对端客户端
func (u *User) ListenMessage() {
	for {
		// 从管道中取出消息
		msg := <-u.C
		// 转换成二进制发送给对端客户端（发送后换行）
		u.conn.Write([]byte(msg + "\n"))
	}
}
