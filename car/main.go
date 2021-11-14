package main

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func main() {
	s := "019230da18a700"
	str := stringMD5(s)

	subStr := strings.ToUpper(str[1:7])
	fmt.Println(subStr)
	var i int
	var j int

	paramString := ""
	for i = 0; i < len(subStr); i = j {
		j = i + 1

		c := subStr[i:j][0]

		stringBuild := strings.Builder{}
		stringBuild.WriteString(paramString)
		tmp := c % 10 * uint8(i+2) % 10

		stringBuild.WriteString(fmt.Sprintf("%d", tmp))
		paramString = stringBuild.String()
	}

	fmt.Println(paramString)
}

func stringMD5(param string) string {
	arrByte := []byte(param)
	arrMD5 := md5.Sum(arrByte)
	i := 0
	stringBuf := strings.Builder{}
	for j := i; j < len(arrMD5); j++ {
		i = int(arrMD5[j] & 0xFF)
		fmt.Println(i)
		if i < 16 {
			stringBuf.WriteString("0")
		}
		stringBuf.WriteString(fmt.Sprintf("%x", i))
	}

	return stringBuf.String()
}
