package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("conn err:", err)
	}
	defer conn.Close()

	go Send(conn)

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn Read err:", err)
			return
		}
		fmt.Println("服务器消息:", string(buf[:n]))
	}
}

func Send(conn net.Conn) {
	strBuf := make([]byte, 1024)

	for {
		n, err := os.Stdin.Read(strBuf)
		if err != nil {
			fmt.Println("os.Stdin.Read err:", err)
		}
		_, err = conn.Write(strBuf[:n])
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return
		}
	}
}
