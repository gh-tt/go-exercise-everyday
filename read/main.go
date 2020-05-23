package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var maxGoChan = make(chan int, 200)

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

func readIp() []string {
	buf, err := ioutil.ReadFile("D:\\www\\go-exercise-everyday\\read\\sjc-ip.txt")
	if err != nil {
		fmt.Println("read file err", err)
		panic(err)
	}

	ips := strings.Split(string(buf), "\n")

	fmt.Println(ips[0])
	s := strings.Split(ips[0], ".")
	fmt.Println(s, len(s))
	s1 := strings.TrimRight(ips[0], ".")
	fmt.Println(s1)

	return ips[:len(ips)-1]
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
	txt, _ := os.OpenFile("D:\\www\\go-exercise-everyday\\read\\sjc-ip2.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
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

func readIp2(ipFileDir string) []string {
	buf, err := ioutil.ReadFile(ipFileDir)
	if err != nil {
		fmt.Println("read file err", err)
		panic(err)
	}

	ips := strings.Split(string(buf), "\n")

	//ips = ips[:len(ips)-1]
	fmt.Println(ips)
	tmp := make([]string, 10)
	ip := strings.Split(ips[0], ".")
	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(254)
		tmp = append(tmp, fmt.Sprintf("%s.%s.%s.%v", ip[0], ip[1], ip[2], num))
	}
	fmt.Println(tmp)
	return ips
}
