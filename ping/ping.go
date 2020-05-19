package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
)

func RunCMD(command string) (string, error) {
	//in := bytes.NewBuffer(nil)
	cmd := exec.Command("/bin/sh", "-c", command)
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

	return string(opBytes), nil
}

func ping(ip string) (string, error) {
	cmd := "ping -c 2 " + ip
	return RunCMD(cmd)
}

var ips []string
var maxGoChan chan int

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	baseIp := "192.168.188."

	for i := 1; i < 8; i++ {
		ips = append(ips, baseIp+strconv.Itoa(i))
	}

	maxGoChan = make(chan int, 5)
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
	//fmt.Println(len(ch))
	/*for i := 0; i < len(ips); i++ {
		fmt.Println(<-ch)
	}*/
}
