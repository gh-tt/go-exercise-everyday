package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

var maxGoChan = make(chan int, 500)

func main() {

	ips := regexpIp()
	ch := make(chan int, len(ips))

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		for i := 0; i < len(ips); i++ {
			<-ch
		}
		wg.Done()
	}()

	for _, ip := range ips {

		maxGoChan <- 1

		go func(ip string) {
			url := fmt.Sprintf("http://%s/cdn-cgi/trace", ip)
			if httpGet(url) {
				//fmt.Println("write...")
				write(ip)
			}
			ch <- 1
		}(ip)

	}
	wg.Wait()
}

func httpGet(url string) bool {
	resp, err := http.Get(url)

	<-maxGoChan

	if err != nil {
		fmt.Println("http get err", err)
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("read body err", err)
		return false
	}
	//fmt.Println(string(body))
	re := strings.Contains(string(body), "colo=SJC")
	return re
}

func write(ip string) {
	txt, _ := os.OpenFile("D:\\www\\go-exercise-everyday\\read\\sjc-ip.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	defer txt.Close()
	txt.WriteString(ip + "\n")
}

func regexpIp() []string {
	buf, err := ioutil.ReadFile("D:\\www\\go-exercise-everyday\\read\\cf-ip.txt")
	if err != nil {
		fmt.Println("read file err", err)
		panic(err)
	}

	reg, err := regexp.Compile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
	if err != nil {
		fmt.Println("regexp compile err", err)
		panic(err)
	}
	ips := reg.FindAllString(string(buf), -1)
	return ips
}
