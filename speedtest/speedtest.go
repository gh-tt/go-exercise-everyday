package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
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
	t := time.Now()
	viper.AddConfigPath("D:\\www\\go-exercise-everyday\\speedtest")
	viper.AddConfigPath("D:\\go-project\\go-exercise-everyday\\speedtest")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	url := viper.GetString("downloadUrl")
	port := viper.GetString("downloadPort")
	ipFileDir := viper.GetString("ipFileDir")

	downloadCount := viper.GetInt("maxDownloadCount")
	minSpeed := viper.GetFloat64("minSpeed")
	ips = readIp(ipFileDir)
	if len(ips) < 1 {
		return
	}
	timeout := setTimeout(url, ips[0], port, minSpeed) * downloadCount
	if timeout > 180 {
		timeout = 180
	}
	fmt.Println("set timeout :", timeout)

	maxGoRoutine := viper.GetInt("maxGoRoutine")
	maxGoChan = make(chan int, maxGoRoutine)
	for i := 0; i < len(ips); i++ {
		maxGoChan <- 1
		wg.Add(1)
		fmt.Println("loop i is", i)
		go Loop(downloadCount, url, ips[i], port, timeout)
		fmt.Println("main i", i)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	}
	wg.Wait()
	resSort()
	modifyDnsAndWrite()
	fmt.Println(time.Since(t))
}

func setTimeout(url, ip, port string, minSpeed float64) int {
	if minSpeed < 1 {
		panic("最小下载速度必须设置大于等于1MB/s")
	}
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	remoteAddr := ip + ":" + port
	length := getRespDataSize(ctx, url, remoteAddr)
	size, _, _ := generateData(length, 1*time.Second)
	if size == 0 {
		panic("获取待测试下载文件大小失败，请重试")
	}
	res := int(size/minSpeed) + 1
	return res
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

func Loop(count int, url, ip, port string, timeout int) {
	defer func() {
		<-maxGoChan
		wg.Done()
	}()
	if count <= 0 {
		return
	}
	remoteAddr := ip + ":" + port

	ch := make(chan int, 1)
	lenChan := make(chan int64, 1)
	timeChan := make(chan time.Duration, 1)
	var l int64
	var dur time.Duration
	ctx, cancel := context.WithCancel(context.Background())
	go func(c context.Context) {
		t := time.Now()
		for i := 0; i < count; i++ {
			select {
			case <-c.Done():
				fmt.Println("loop exit", i)
				return
			default:
				l += getRespDataSize(c, url, remoteAddr)
			}
			fmt.Println("********************************")
		}
		dur := time.Since(t)
		lenChan <- l
		timeChan <- dur
		ch <- 1
	}(ctx)
	select {
	case <-ch:
		l = <-lenChan
		dur = <-timeChan
		fmt.Println(l, dur, "size and dur")
		fmt.Println("-----------------------------------------------")
	case <-time.After(time.Duration(timeout) * time.Second):
		fmt.Println("time out")
		cancel()
		return
	}
	size, useTime, speed := generateData(l, dur)

	mu.Lock()
	speedStatSlice = append(speedStatSlice, SpeedStat{Ip: ip, Size: size, UseTime: useTime, Speed: speed})
	mu.Unlock()
}

func generateData(l int64, dur time.Duration) (float64, float64, float64) {
	size := float64(l) / (1024 * 1024)
	str := fmt.Sprintf("%.2f", size)
	size, _ = strconv.ParseFloat(str, 64)

	useTime := float64(dur) / 1e9
	str = fmt.Sprintf("%.2f", useTime)
	useTime, _ = strconv.ParseFloat(str, 64)

	speed := float64(l) / (1024 * 1024) * 1e9 / float64(dur)
	str = fmt.Sprintf("%.2f", speed)
	speed, _ = strconv.ParseFloat(str, 64)
	return size, useTime, speed
}

func getRespDataSize(ctx context.Context, url, remoteAddr string) int64 {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("req err: ", err)
		return 0
	}
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
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("resp err: ", err)
		return 0
	}
	defer resp.Body.Close()
	var length int64
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(" read resp err :", err)
			return 0
		}
		length += int64(n)
	}
	/*body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(" read resp err :", err)
		return 0
	}
	fmt.Println(len(body))
	return int64(len(body))*/
	fmt.Println(length)
	return length
}

//url 请求的URL,例如：https://www.baidu.com/123/img.png
// remoteAddr 请求的远程服务器地址，例如:192.168.1.1:443
func httpGet(ctx context.Context, url, remoteAddr string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(c context.Context, network, addr string) (net.Conn, error) {
				add, _ := net.ResolveTCPAddr("tcp", remoteAddr)
				conn, err := net.DialTCP("tcp", nil, add)
				return conn, err
			},
		},
	}
	req.Header.Set("User-Agent", "go-http")
	return client.Do(req)
}

func modifyDnsAndWrite() {
	if len(speedStatSlice) == 0 {
		return
	}
	if !viper.GetBool("dns.modifyEnable") {
		fmt.Println("do not need modify dns")
		return
	}
	if len(speedStatSlice) > 0 {
		ip := speedStatSlice[0].Ip

		data := make(url.Values)
		data["login_token"] = []string{viper.GetString("dns.dnspodToken")}
		data["domain"] = []string{viper.GetString("dns.domain")}
		data["sub_domain"] = []string{viper.GetString("dns.subDomain")}
		data["record_id"] = []string{viper.GetString("dns.recordId")}
		data["record_type"] = []string{viper.GetString("dns.recordType")}
		data["record_line"] = []string{viper.GetString("dns.recordLine")}
		data["value"] = []string{ip}

		if speedStatSlice[0].Speed >= viper.GetFloat64("dns.speedLimit") {
			_, _ = http.PostForm("https://dnsapi.cn/Record.Modify", data)
			fmt.Println("modifyDns success")
		} else {
			fmt.Println("ip 不符合要求")
		}
	}
	betterIp := speedStatSlice
	if len(betterIp) > 5 {
		betterIp = betterIp[:5]
	}
	write(betterIp)
}

func write(betteIp []SpeedStat) {
	osType := runtime.GOOS
	if osType == "windows" {
		fmt.Println("windows不需要写日志")
		return
	}
	path := "/data/gotools/speedtest/result.txt"
	txt, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	defer txt.Close()
	t := time.Now()
	str := fmt.Sprintln(betteIp, "   ---   ", t.Format("2006-01-02 15:04:05"))
	n, err := txt.WriteString(str)
	if n != len(str) {
		panic(err)
	}
}
