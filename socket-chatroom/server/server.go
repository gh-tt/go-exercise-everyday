package main

import (
	"fmt"
	"net"
	_ "net/http/pprof"
)

type User struct {
	Addr string
}

var clients = make(map[net.Conn]User)

var messages = make(chan string, 100)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	fmt.Println("开始监听,等待用户连接")

	defer listener.Close()

	go broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		go HandleConn(conn)
	}
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()

	fmt.Println(addr, "已连接")
	user := User{Addr: addr}
	clients[conn] = user

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err", err)
			delete(clients, conn)
			return
		}

		msg := addr + ":" + string(buf[:n])
		messages <- msg
	}
}

func broadcast() {
	for msg := range messages {
		msg := msg
		go func() {
			for cli := range clients {
				cli.Write([]byte(msg))
			}
		}()
	}
}
