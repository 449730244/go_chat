package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

//声明一个Sever对象
type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User //保存当前在线用户
	mapLock   sync.RWMutex
	Message   chan string //用户消息
}

//实例化一个Server对象
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//启动服务
func (this *Server) Start() {
	//监听服务
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("Server Start err", err)
		return
	}
	//执行结束后关闭连接
	defer listener.Close()

	//启动消息监听
	go this.ListenMessage()
	//接收客户端连接
	for {
		//accept
		conn, err := listener.Accept()
		//如果连接失败跳出当前连接
		if err != nil {
			fmt.Println("Client Connent Fail", err)
			continue
		}
		//do handel
		go this.Handel(conn)
	}
}

//监听Message广播消息,一旦有消息发送给所有在线用户
func (this *Server) ListenMessage() {
	for {
		//获取消息
		msg := <-this.Message
		this.mapLock.Lock()
		//获取在线用户
		for _, cli := range this.OnlineMap {
			cli.C <- msg //将消息发送给用户
		}
		this.mapLock.Unlock()
	}
}

//广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg //消息拼接
	this.Message <- sendMsg                                  //讲消息放入channel中
}

//具体业务操作
func (this *Server) Handel(conn net.Conn) {
	//获取user
	user := NewUser(conn, this)
	//用户上线
	user.Online()
	go func() {
		buf := make([]byte, 4896)
		for {
			n, err := conn.Read(buf) //读取消息
			//如果读取为0
			if n == 0 {
				user.Offline() //用下线
				return
			}
			//错误判断
			if err != nil && err != io.EOF {
				fmt.Println("conn Read err", err)
				return
			}
			//提取用户消息，去掉\n
			msg := string(buf[:n-1])
			//广播用户消息
			user.DoMessage(msg)
		}
	}()

}
