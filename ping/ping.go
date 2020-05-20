package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"sync"
)

const (
	maxGOROUTINE = 255
	maxPingNum   = 100
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
	baseIp := "104.18.0."
	runtime.GOMAXPROCS(runtime.NumCPU())

	for i := 1; i < 255; i++ {
		ips = append(ips, baseIp+strconv.Itoa(i))
	}

	maxGoChan = make(chan int, maxGOROUTINE)
	ch := make(chan string, len(ips))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < len(ips); i++ {
			fmt.Println(<-ch)
		}
		wg.Done()
	}()

	for i := 0; i < len(ips); i++ {

		maxGoChan <- 1

		go func(i int) {
			pingRe, _ := ping(ips[i])
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
	//fmt.Println(ipLossPercentMap)
	//fmt.Println(len(ch))
	/*for i := 0; i < len(ips); i++ {
		fmt.Println(<-ch)
	}*/
}

func write(betteIp []PingStat) {
	txt, _ := os.OpenFile("ping.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer txt.Close()

	str := fmt.Sprintln(betteIp)
	n, err := txt.WriteString(str)
	if n != len(str) {
		panic(err)
	}
}

func ping(ip string) (string, error) {
	return RunCMD(ip, maxPingNum)
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
