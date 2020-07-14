package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

const defaultTcpPort = 443
const tcpConnectTimeout = time.Millisecond * 500
const failTime = 4

type CloudflareIPData struct {
	ip            string
	pingCount     int
	pingTime      float64
	pingReceived  int
	recvRate      float64
	downloadSpeed float64
	downloadTime  int
}

func (cf *CloudflareIPData) getRecvRate() float64 {
	if cf.recvRate == 0 {
		cf.recvRate = float64(cf.pingReceived) / float64(cf.pingCount) * 100
	}
	return cf.recvRate
}

func ExportTxt(filepath string, data []CloudflareIPData) {
	txt, _ := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	defer txt.Close()
	t := time.Now()
	str := fmt.Sprintln(data, "   ---   ", t.Format("2006-01-02 15:04:05"))
	n, err := txt.WriteString(str)
	if n != len(str) {
		panic(err)
	}
}

func loadIp() []string {
	buf, err := ioutil.ReadFile(Conf.ipFilePath)
	if err != nil {
		fmt.Println("read ip file err", err)
		panic(err)
	}

	ips := strings.Split(string(buf), "\n")
	ips = ips[:len(ips)-1]
	ipList := make([]string, 0)

	count := Conf.selectCountEveryIp
	if count <= 0 || count > 255 {
		panic("每个ip段选择的ip数量,不能为0且小于等于255")
	}

	rand.Seed(time.Now().UnixNano())
	for _, v := range ips {
		ip := strings.Split(v, ".")
		for i := 0; i < count; i++ {
			num := rand.Intn(254) + 1
			ipList = append(ipList, fmt.Sprintf("%s.%s.%s.%v", ip[0], ip[1], ip[2], num))
		}
	}

	return ipList
}
