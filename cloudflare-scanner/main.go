package main

import (
	"fmt"
	"sort"
	"sync"
)

func main() {
	initConfig()
	var wg sync.WaitGroup
	var mu sync.Mutex
	var data = make([]CloudflareIPData, 0)
	ips := loadIp()
	pingRoutine := make(chan bool, Conf.pingRoutine)
	for _, ip := range ips {
		wg.Add(1)
		pingRoutine <- false
		go pingGoroutine(&wg, &mu, ip, Conf.pingCount, &data, pingRoutine)
	}
	wg.Wait()
	sort.Slice(data, func(i, j int) bool {
		if data[i].getRecvRate() != data[j].getRecvRate() {
			return data[i].getRecvRate() > data[j].getRecvRate()
		}
		return data[i].pingTime < data[j].pingTime
	})
	fmt.Println(data)
}
