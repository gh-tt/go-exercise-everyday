package main

import (
	"bufio"
	"context"
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

var maxGoChan = make(chan int, 2)
var fail int
var failMu sync.Mutex

func main() {
	mu := sync.Mutex{}
	ips := regexpIp()
	wg := sync.WaitGroup{}
	fmt.Println(ips)
	fmt.Println(len(ips))

	for _, ip := range ips {

		maxGoChan <- 1
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			url := fmt.Sprintf("http://%s/cdn-cgi/trace", ip)
			if httpGet(url) {
				defer mu.Unlock()
				mu.Lock()
				write(ip)
			}
		}(ip)

	}

	wg.Wait()
	fmt.Println(fail)
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
	defer func() {
		<-maxGoChan
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("req err: ", err)
		failMu.Lock()
		fail++
		failMu.Unlock()
		return false
	}
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("resp err: ", err)
		failMu.Lock()
		fail++
		failMu.Unlock()
		return false
	}
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("read body err:", err)
			failMu.Lock()
			fail++
			failMu.Unlock()
			return false
		}
		//fmt.Println(string(body))
		re := strings.Contains(string(body), "colo=SJC")
		return re
	} else {
		failMu.Lock()
		fail++
		failMu.Unlock()
		return false
	}

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

func regexpIp2() {
	fi, _ := os.Open("D:\\www\\go-exercise-everyday\\read\\cf-ip.txt")
	s := make([]string, 0)
	defer fi.Close()
	reader := bufio.NewReader(fi)
	for {
		tmp, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		//fmt.Println(tmp)
		s = append(s, tmp)
	}
	newS := make([]string, 0)
	for i := 0; i < len(s)-1; i++ {
		if s[i+1] == "colo=SJC\n" {
			newS = append(newS, s[i])
		}
	}
	ips := make([]string, 0)
	reg, _ := regexp.Compile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
	for _, v := range newS {
		tmp := reg.FindAllString(v, -1)
		ips = append(ips, tmp...)
	}
	for _, v := range ips {
		write(v)
	}
	//fmt.Println(ips, len(ips))
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
