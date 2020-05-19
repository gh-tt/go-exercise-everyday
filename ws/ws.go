package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"

	"time"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("hello"))
	var (
		wsConn *websocket.Conn
		err    error
		conn   *Connection
		data   []byte
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = InitConnection(wsConn); err != nil {
		goto ERR
	}

	// 启动线程，不断发消息
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(100 * time.Second)
		}
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}
func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("D:\\www\\go-exercise-everyday\\ws\\index.html")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	t.Execute(w, nil)
}
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":7777", nil)
}
