package main

import (
	"fmt"
	"net"
)

//声明一个Sever对象
type Server struct {
	Ip   string
	Port int
}

//实例化一个Server对象
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
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

//具体业务操作
func (this *Server) Handel(conn net.Conn) {
	fmt.Println("Server connent success\n")
}
