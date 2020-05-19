package main

import (
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"os/exec"
)

func main() {
	cmd := exec.Command("ping", "www.baidu.com")
	outinfo := bytes.Buffer{}
	cmd.Stdout = &outinfo

	err := cmd.Start()
	if err != nil {
		fmt.Println("cmd start err:", err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("cmd wait err:", err)
	} else {
		dec := mahonia.NewDecoder("gbk")

		_, newOut, err := dec.Translate(outinfo.Bytes(), true)
		if err != nil {
			fmt.Println("gbk translate err", err)
		} else {
			fmt.Println(string(newOut))
		}

	}
}
