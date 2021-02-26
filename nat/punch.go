package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	serverPort, serverAddr string
)

func main() {
	flag.StringVar(&serverPort, "s", "", "服务端模式，监听的端口")
	flag.StringVar(&serverAddr, "c", "", "客户端模式，要连接的服务器地址")
	flag.Parse()

	if serverPort != "" {
		serverRun()
		return
	}
	if serverAddr != "" {
		clientRun()
		return
	}
	log.Fatal("请指定运行模式（服务端 -s listenPort |客户端 -c serverIP:port）")
}

func clientRun() {

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Println("listen udp err:", err)
		return
	}

	raddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("resolve udp raddr err:", err)
		return
	}

	sendData("1", conn, raddr)

	var coupleAddr *net.UDPAddr
	for {
		data, addr := receiveData(conn)
		fmt.Println("remote addr:", addr.String())
		if addr.String() == raddr.String() {
			coupleAddr, err = net.ResolveUDPAddr("udp", string(data))
			if err != nil {
				log.Fatal("couple addr err:", err)
			}
			break
		}
	}

	canSend := false
	for !canSend {
		sendData("2", conn, coupleAddr)
		for {
			data, addr := receiveData(conn)
			fmt.Println("remote addr:", addr.String())
			if addr.String() == coupleAddr.String() {
				fmt.Println(data)
				canSend = true
				break
			}
		}
	}

	go func() {
		for {
			data, addr := receiveData(conn)
			if addr.String() == coupleAddr.String() {
				fmt.Println(data)
			}
		}
	}()

	go func() {
		for {
			sendData("hello", conn, coupleAddr)
			time.Sleep(time.Second)
		}
	}()

	select {}
}

func serverRun() {
	addr := fmt.Sprintf(":%s", serverPort)
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("laddr resolve err :", err)
		return
	}
	listenConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Println("listen udp err :", err)
		return
	}
	defer listenConn.Close()

	for {
		_, firstAddr := receiveData(listenConn)
		_, secondAddr := receiveData(listenConn)

		fmt.Printf("%s <=> %s\n", firstAddr.String(), secondAddr.String())

		sendData(firstAddr.String(), listenConn, secondAddr)
		sendData(secondAddr.String(), listenConn, firstAddr)
	}

}

func sendData(data string, conn *net.UDPConn, toAddr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte(data), toAddr)
	if err != nil {
		log.Fatal("write data err:", err)
		return
	}
	fmt.Println("send data success: ", data)

}

func receiveData(conn *net.UDPConn) (data []byte, addr *net.UDPAddr) {
	buf := make([]byte, 50)

	n, sendAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Fatal("read udp err:", err)
	}
	fmt.Println("receive data success: ", string(buf[:n]))
	return buf[:n], sendAddr
}
