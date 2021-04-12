package main

func main() {
	//创建服务
	server := NewServer("127.0.0.1", 8989)
	server.Start()
}
