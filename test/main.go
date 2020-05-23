package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func main() {
	readIp2("D:\\www\\go-exercise-everyday\\read\\sjc-ip.txt")
}

func readIp2(ipFileDir string) []string {
	buf, err := ioutil.ReadFile(ipFileDir)
	if err != nil {
		fmt.Println("read file err", err)
		panic(err)
	}

	ips := strings.Split(string(buf), "\n")

	ips = ips[:len(ips)-1]

	ipList := make([]string, 0)

	rand.Seed(time.Now().UnixNano())
	for _, v := range ips {
		ip := strings.Split(v, ".")
		for i := 0; i < 10; i++ {
			num := rand.Intn(255)
			if num != 0 {
				ipList = append(ipList, fmt.Sprintf("%s.%s.%s.%v", ip[0], ip[1], ip[2], num))
			}
		}
	}

	return ipList
}
