package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	h  bool
	u  string
	ip string
	t  int
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&u, "u", "https://speed.cloudflare.com/__down?bytes=10000000", "cf download url")
	flag.StringVar(&ip, "ip", "172.64.195.171", "download test use ip")
	flag.IntVar(&t, "t", 10, "download test time")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `cf single ip speedtest
Usage: speedtest [-h] [-u url] [-t time]

Options:
	`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	fmt.Println("*****开始测速*****")
	fmt.Println("测速资源url：",u)
	fmt.Println("测速ip:",ip,"     ","测速时间:",t,"秒")
	speedGoRoutine(u, ip, t)
}

func speedGoRoutine(url, ip string, downSecond int) {
	d := time.Now().Add(time.Duration(downSecond) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	var downloadDataSize int64

Loop:
	for {
		select {
		case <-ctx.Done():
			break Loop
		default:
			DownloadHandler(ctx, url, ip, &downloadDataSize)
		}
	}
	speed := float64(downloadDataSize) / 1024 / 1024 / float64(downSecond)
	fmt.Printf("当前ip速度测试结果为:%.2fMB/s\n",speed)
}

func DownloadHandler(ctx context.Context, url, ip string, downloadDataSize *int64) (bool, int64) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, 0
	}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(c context.Context, network, addr string) (net.Conn, error) {
				conn, err := (&net.Dialer{}).DialContext(c, network, ip+":443")
				//conn, err := net.Dial("tcp", ip+":443")
				return conn, err
			},
			DisableKeepAlives: true,
		},
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		buf := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buf)
			*downloadDataSize += int64(n)
			if err != nil && err == io.EOF {
				break
			} else if err != nil {
				return false, 0
			}

		}
		return true, 0
	} else {
		return false, 0
	}
}
