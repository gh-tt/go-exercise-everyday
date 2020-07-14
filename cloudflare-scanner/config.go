package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	selectCountEveryIp int
	ipFilePath         string
	pingRoutine        int
	pingCount          int
	speedTestCount     int
	downloadSecond     int
	downloadRoutine    int
}

var Conf *Config

func initConfig() {
	viper.AddConfigPath("./")
	viper.AddConfigPath("/data/gotools/cloudflare-scanner/")
	viper.AddConfigPath("D:\\www\\go-exercise-everyday\\cloudflare-scanner")
	viper.AddConfigPath("D:\\go-project\\go-exercise-everyday\\cloudflare-scanner")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件失败，请检查 config.yaml 配置文件是否存在: %v", err))
	}
	Conf = newConfig()
}
func newConfig() *Config {
	pingRoutine := viper.GetInt("pingRoutine")
	if pingRoutine <= 0 {
		pingRoutine = 100
	}
	pingCount := viper.GetInt("pingCount")
	if pingCount <= 0 {
		pingCount = 10
	}
	speedTestCount := viper.GetInt("downloadTestCount")
	if speedTestCount <= 0 {
		speedTestCount = 10
	}
	downloadSecond := viper.GetInt("downloadSecond")
	if downloadSecond <= 0 {
		downloadSecond = 10
	}
	downloadRoutine := viper.GetInt("downloadRoutine")
	if downloadRoutine <= 0 {
		downloadRoutine = 1
	}
	return &Config{
		selectCountEveryIp: viper.GetInt("selectCountEveryIp"),
		ipFilePath:         viper.GetString("ipFileDir"),
		pingRoutine:        pingRoutine,
		pingCount:          pingCount,
		speedTestCount:     speedTestCount,
		downloadSecond:     downloadSecond,
		downloadRoutine:    downloadRoutine,
	}
}
