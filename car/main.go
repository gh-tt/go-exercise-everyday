package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func main() {
	router := gin.Default()

	router.GET("/get-car-pwd", getCarPwd)

	err := router.Run(":21111")
	log.Fatalln(err)
	//deviceId := "019230da18a700"

}

func getCarPwd(c *gin.Context) {
	deviceId := c.Query("id")
	if deviceId == "" {
		c.String(http.StatusNotFound, "404 page not found")
		return
	}
	str := stringMD5(deviceId)
	pwd := generatePwd(str)
	fmt.Println("deviceID:", deviceId, "; password:", pwd)

	c.String(http.StatusOK, pwd)
	return
}

func generatePwd(md5String string) string {
	subStr := strings.ToUpper(md5String[1:7])
	//fmt.Println(subStr)
	var i int
	var j int

	password := ""
	for i = 0; i < len(subStr); i = j {
		j = i + 1

		c := subStr[i:j][0]

		stringBuild := strings.Builder{}
		stringBuild.WriteString(password)
		tmp := c % 10 * uint8(i+2) % 10

		stringBuild.WriteString(fmt.Sprintf("%d", tmp))
		password = stringBuild.String()
	}

	//fmt.Println(password)
	return password
}

func stringMD5(param string) string {
	arrByte := []byte(param)
	arrMD5 := md5.Sum(arrByte)
	i := 0
	stringBuf := strings.Builder{}
	for j := i; j < len(arrMD5); j++ {
		i = int(arrMD5[j] & 0xFF)
		//fmt.Println(i)
		if i < 16 {
			stringBuf.WriteString("0")
		}
		stringBuf.WriteString(fmt.Sprintf("%x", i))
	}

	return stringBuf.String()
}
