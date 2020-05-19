package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"sync"
)

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
	reg := regexp.MustCompile("(\\d+)% packet loss")
	result := reg.Find(opBytes)
	fmt.Println(string(result[:len(result)-13]))

	lossPercent := result[:len(result)-13]
	ipLossPercentMap[ip], _ = strconv.Atoi(string(lossPercent))
	return string(opBytes), nil
}

func ping(ip string) (string, error) {
	return RunCMD(ip, 100)
}

var ips []string
var maxGoChan chan int
var ipLossPercentMap map[string]int

func main() {
	ipLossPercentMap = make(map[string]int)
	runtime.GOMAXPROCS(runtime.NumCPU())

	baseIp := "104.18.0."

	for i := 1; i < 254; i++ {
		ips = append(ips, baseIp+strconv.Itoa(i))
	}

	maxGoChan = make(chan int, 10)
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

	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range ipLossPercentMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})
	fmt.Println(ss[:5])
	//fmt.Println(ipLossPercentMap)
	//fmt.Println(len(ch))
	/*for i := 0; i < len(ips); i++ {
		fmt.Println(<-ch)
	}*/
}
