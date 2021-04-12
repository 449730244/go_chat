package main

import "net"

//声明一个user对象
type User struct {
	Name   string
	Addr   string
	C      chan string //一个用户一个channel
	conn   net.Conn
	server *Server
}

//创建一个user对象
func NewUser(conn net.Conn, server *Server) *User {
	//获取当前用户IP地址
	addr := conn.RemoteAddr().String()
	user := &User{
		Name:   addr,
		Addr:   addr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动消息监听
	go user.ListenMsg()
	return user
}

//监听当前user的channel中的消息
func (this *User) ListenMsg() {
	for {
		msg := <-this.C                     //去除channle中的信息
		this.conn.Write([]byte(msg + "\n")) //将消息转成byte发送
	}
}

//用户上线
func (this *User) Online() {
	//用户上线加入到OnlineMap中
	this.server.mapLock.Lock()              //加锁
	this.server.OnlineMap[this.Name] = this //讲用户加入到OnlineMap
	this.server.mapLock.Unlock()            //释放锁
	//广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}

//用户下线
func (this *User) Offline() {
	this.server.mapLock.Lock()               //加锁
	delete(this.server.OnlineMap, this.Name) //将用户从OnlineMap中删除
	this.server.mapLock.Unlock()             //释放锁
	this.server.BroadCast(this, "已下线")
}

//用户消息处理
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}
