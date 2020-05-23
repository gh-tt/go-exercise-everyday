package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ips           []string
	maxGoChan     chan int
	pingStatSlice []PingStat
)

type PingStat struct {
	Ip   string  //ip
	Loss int     //ping packet loss per
	Rtt  float64 //ping rtt avg
}

func main() {
	viper.AddConfigPath("./")
	viper.AddConfigPath("/data/gotools/cf-ping/")
	//viper.AddConfigPath("/data/go/go-exercise-everyday/ping")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	ipFileDir := viper.GetString("ipFileDir")
	ips = readIp(ipFileDir)

	maxGoRoutine := viper.GetInt("maxGoRoutine")
	maxGoChan = make(chan int, maxGoRoutine)
	ch := make(chan string, len(ips))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < len(ips); i++ {
			fmt.Println(<-ch)
		}
		wg.Done()
	}()

	maxPingCount := viper.GetInt("maxPingCount")

	for i := 0; i < len(ips); i++ {

		maxGoChan <- 1

		go func(i int) {
			pingRe, _ := ping(ips[i], maxPingCount)
			if pingRe == "" {
				fmt.Println("ping fail")
			}
			ch <- pingRe

		}(i)
	}

	wg.Wait()

	sort.Slice(pingStatSlice, func(i, j int) bool {
		return pingStatSlice[i].Loss < pingStatSlice[j].Loss
	})

	var betterIp []PingStat
	if len(pingStatSlice) > 0 {
		for _, v := range pingStatSlice {
			if v.Loss == pingStatSlice[0].Loss {
				betterIp = append(betterIp, v)
			} else {
				break
			}
		}
	}

	sort.Slice(betterIp, func(i, j int) bool {
		return betterIp[i].Rtt < betterIp[j].Rtt
	})
	fmt.Println(pingStatSlice)
	fmt.Println("--------------------------------------------------")
	fmt.Println(betterIp)
	if len(betterIp) > 5 {
		betterIp = betterIp[:5]
	}
	write(betterIp)
	modifyDns(betterIp)
}

func modifyDns(betterIp []PingStat) {
	if !viper.GetBool("dns.modifyEnable") {
		fmt.Println("do not need modify dns")
		return
	}

	if len(betterIp) > 0 {
		ip := betterIp[0].Ip

		data := make(url.Values)
		data["login_token"] = []string{viper.GetString("dns.dnspodToken")}
		data["domain"] = []string{viper.GetString("dns.domain")}
		data["sub_domain"] = []string{viper.GetString("dns.subDomain")}
		data["record_id"] = []string{viper.GetString("dns.recordId")}
		data["record_type"] = []string{viper.GetString("dns.recordType")}
		data["record_line"] = []string{viper.GetString("dns.recordLine")}
		data["value"] = []string{ip}

		if betterIp[0].Loss <= viper.GetInt("dns.lossLimit") && betterIp[0].Rtt <= viper.GetFloat64("dns.rttLimit") {
			_, _ = http.PostForm("https://dnsapi.cn/Record.Modify", data)
			fmt.Println("modifyDns success")
		} else {
			fmt.Println("ip 不符合要求")
		}

	}
}

func write(betteIp []PingStat) {
	txt, _ := os.OpenFile("/data/gotools/cf-ping/ping.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	defer txt.Close()

	str := fmt.Sprintln(betteIp)
	n, err := txt.WriteString(str)
	if n != len(str) {
		panic(err)
	}
}

func ping(ip string, count int) (string, error) {
	return RunCMD(ip, count)
}

func RunCMD(ip string, count int) (string, error) {
	//in := bytes.NewBuffer(nil)
	cmd := exec.Command("/bin/sh", "-c", "ping -c "+strconv.Itoa(count)+" "+ip)
	//cmd.Stdin = in
	//in.WriteString(command + "\n")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("stdout err", err)
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("cmd start err", err)
		log.Fatal(err)
	}

	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("read err", err)
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("cmd Wait err: ", err)
		//log.Fatal(err)
	}

	<-maxGoChan
	reg1 := regexp.MustCompile("(\\d+)% packet loss")
	result := reg1.FindSubmatch(opBytes)
	fmt.Println("loss per:", string(result[1]))
	lossPercent := result[1]
	pLp, _ := strconv.Atoi(string(lossPercent))

	reg2, _ := regexp.Compile(`/([0-9]+\.?[0-9]+)/`)
	res2 := reg2.FindSubmatch(opBytes)
	if res2 != nil {
		rtt, _ := strconv.ParseFloat(string(res2[1]), 64)
		pingStatSlice = append(pingStatSlice, PingStat{ip, pLp, rtt})
	}

	return string(opBytes), nil
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
