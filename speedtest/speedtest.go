package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ips            []string
	maxGoChan      chan int
	speedStatSlice []SpeedStat
	mu             sync.Mutex
	wg             sync.WaitGroup
)

type SpeedStat struct {
	Ip      string  //ip
	Size    float64 //数据大小，单位MB
	UseTime float64 //总用时 单位s
	Speed   float64 //速度 单位MB/s
}

func main() {
	viper.AddConfigPath("D:\\www\\go-exercise-everyday\\speedtest")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	ipFileDir := viper.GetString("ipFileDir")
	ips = readIp(ipFileDir)
	url := "https://storage.idx0.workers.dev/Images/public-notion-06b4a73f-0d4e-4b8f-b273-77becf84a0b3.png"
	port := "443"

	maxGoRoutine := viper.GetInt("maxGoRoutine")
	maxGoChan = make(chan int, maxGoRoutine)
	for i := 0; i < len(ips); i++ {
		maxGoChan <- 1
		wg.Add(1)
		go Loop(1, url, ips[i], port)
	}
	wg.Wait()
	resSort()
}
func resSort() {
	sort.Slice(speedStatSlice, func(i, j int) bool {
		return speedStatSlice[i].Speed > speedStatSlice[j].Speed
	})
	if len(speedStatSlice) > 5 {
		fmt.Println(speedStatSlice[:5])
	} else {
		fmt.Println(speedStatSlice)
	}
}

func readIp(ipFileDir string) []string {
	buf, err := ioutil.ReadFile(ipFileDir)
	if err != nil {
		fmt.Println("read ip file err", err)
		panic(err)
	}

	ips := strings.Split(string(buf), "\n")
	ips = ips[:len(ips)-1]
	ipList := make([]string, 0)

	count := viper.GetInt("selectCountEveryIp")
	if count == 0 || count > 255 {
		panic("每个ip段选择的ip数量,不能为0且小于等于255")
	}

	rand.Seed(time.Now().UnixNano())
	for _, v := range ips {
		ip := strings.Split(v, ".")
		for i := 0; i < count; i++ {
			num := rand.Intn(255)
			if num != 0 {
				ipList = append(ipList, fmt.Sprintf("%s.%s.%s.%v", ip[0], ip[1], ip[2], num))
			}
		}
	}

	return ipList
}

func Loop(count int, url, ip, port string) {
	defer func() {
		<-maxGoChan
		wg.Done()
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	remoteAddr := ip + ":" + port

	ch := make(chan int,1)
	lenChan := make(chan int64,1)
	timeChan := make(chan time.Duration,1)
	var l int64
	var dur time.Duration
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				return
			}
		}()
		t := time.Now()
		for i := 0; i < count; i++ {
			l += getRespDataSize(url, remoteAddr)
		}
		dur := time.Since(t)
		lenChan <- l
		timeChan <- dur
		ch <- 1
	}()
	select {
	case <-ch:
		l = <-lenChan
		dur = <-timeChan
		fmt.Println(l,dur,"size and dur")
	case <-time.After(10 * time.Second):
		fmt.Println("time out")
		return
	}

	size := float64(l) / (1024 * 1024)
	str := fmt.Sprintf("%.2f", size)
	size, _ = strconv.ParseFloat(str, 64)

	useTime := float64(dur) / 1e9
	str = fmt.Sprintf("%.2f", useTime)
	useTime, _ = strconv.ParseFloat(str, 64)

	speed := float64(l) / (1024 * 1024) * 1e9 / float64(dur)
	str = fmt.Sprintf("%.2f", speed)
	speed, _ = strconv.ParseFloat(str, 64)

	mu.Lock()
	speedStatSlice = append(speedStatSlice, SpeedStat{Ip: ip, Size: size, UseTime: useTime, Speed: speed})
	mu.Unlock()
}

func getRespDataSize(url, remoteAddr string) int64 {

	resp, err := httpGet(url, remoteAddr)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0
	}
	return int64(len(body))
}

//url 请求的URL,例如：https://www.baidu.com/123/img.png
// remoteAddr 请求的远程服务器地址，例如:192.168.1.1:443
func httpGet(url, remoteAddr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(c context.Context, network, addr string) (net.Conn, error) {
				add, _ := net.ResolveTCPAddr("tcp", remoteAddr)
				conn, err := net.DialTCP("tcp", nil, add)
				return conn, err
			},
		},
	}
	req.Header.Set("User-Agent", "golang-client")
	return client.Do(req)
}
