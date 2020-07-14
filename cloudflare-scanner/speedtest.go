package main

import (
	"net"
	"strconv"
	"sync"
	"time"
)

//bool connectionSucceed float32 time
func ping(ip string) (bool, float64) {
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", ip+":"+strconv.Itoa(defaultTcpPort), tcpConnectTimeout)
	if err != nil {
		return false, 0
	} else {
		var endTime = time.Since(startTime)
		var duration = float64(endTime.Microseconds()) / 1000.0
		_ = conn.Close()
		return true, duration
	}
}

//pingReceived pingTotalTime
func checkConnection(ip string) (int, float64) {
	pingRecv := 0
	var pingTime float64 = 0.0
	for i := 1; i <= failTime; i++ {
		pingSucceed, pingTimeCurrent := ping(ip)
		if pingSucceed {
			pingRecv++
			pingTime += pingTimeCurrent
		}
	}
	return pingRecv, pingTime
}

//return Success packetRecv averagePingTime specificIPAddr
func pingHandler(ip string, pingCount int) (bool, int, float64, string) {
	ipCanConnect := false
	pingRecv := 0
	var pingTime float64 = 0.0

	pingRecvCurrent, pingTimeCurrent := checkConnection(ip)
	if pingRecvCurrent != 0 {
		ipCanConnect = true
		pingRecv = pingRecvCurrent
		pingTime = pingTimeCurrent
	}

	if ipCanConnect {
		for i := failTime; i < pingCount; i++ {
			pingSuccess, pingTimeCurrent := ping(ip)
			if pingSuccess {
				pingRecv++
				pingTime += pingTimeCurrent
			}
		}
		return true, pingRecv, pingTime / float64(pingRecv), ip
	} else {
		return false, 0, 0, ""
	}
}

func pingGoroutine(wg *sync.WaitGroup, mutex *sync.Mutex, ip string, pingCount int, data *[]CloudflareIPData, pingRoutine chan bool) {
	defer func() {
		<-pingRoutine
		wg.Done()
	}()

	success, pingRecv, pingTimeAvg, currentIP := pingHandler(ip, pingCount)
	if success {
		mutex.Lock()
		var cfdata CloudflareIPData
		cfdata.ip = currentIP
		cfdata.pingReceived = pingRecv
		cfdata.pingTime = pingTimeAvg
		cfdata.pingCount = pingCount
		*data = append(*data, cfdata)
		mutex.Unlock()
	}
}

//bool : can download,float32 downloadSpeed
func DownloadSpeedHandler(ip string) (bool, float64) {
	return false, 0
}
