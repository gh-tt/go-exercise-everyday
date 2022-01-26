package main

import (
	"embed"
	"fmt"
	"github.com/lxn/win"
	"gocv.io/x/gocv"
	"os"
	"sync"
	"time"
)

//go:embed static
var fs embed.FS

var server = &Server{androidMap: make(map[string]*Android), mW: &MyMainWindow{}}

var DeviceList []string

var jingjiTimes, maoxianTimes = 0, 0

type Server struct {
	sync.RWMutex
	sync.WaitGroup
	closeCh      chan int
	androidMap   map[string]*Android
	emulatorList []*Emulator
	mW           *MyMainWindow
}

func init() {
	b, err := os.ReadFile("emulator-path.txt")
	if err == nil {
		EmulatorPath = string(b)
	}
}

func main() {
	server.newGui()
}
func (s *Server) startTask() {
	jingjiTimes, maoxianTimes = 0, 0
	if s.mW.JingjiCheck.Checked() {
		jingjiTimes = int(s.mW.JingjiNum.Value())
	}
	if s.mW.MaoXianCheck.Checked() {
		maoxianTimes = int(s.mW.MaoXianNum.Value())
	}
	if s.mW.MaoXianCheck.Checked() && s.mW.MaoXianEnd.Checked() {
		maoxianTimes = -1
	}
	s.closeCh = make(chan int)
	for k := range s.emulatorList {
		//server.androidMap[k].AdbShellInputSwipeL(120,460,120,260,5000)
		//return
		if s.emulatorList[k].checked {
			s.Add(1)
			go s.start(s.emulatorList[k])
		}
	}
}
func (s *Server) start(emulator *Emulator) {
	defer s.Done()
	emulator.Status = "进行中"
	emulator.MaoxianCount = 0
	emulator.JingjiCount = 0
	fmt.Println("config ：", emulator.MaoxianCount, emulator.JingjiCount, maoxianTimes, jingjiTimes)
	emulator.StartApp(APPPackageName + "/" + ClassName)

	advTem := loadTemplate("static/adventure.png")
	startTem := loadTemplate("static/start.png")
	energyTem := loadTemplate("static/energyless.png")
	fightTem := loadTemplate("static/fight.png")

	emptyCardTimes := 0
	shotFail := 0
	for {
		select {
		case <-s.closeCh:
			goto Out
		default:
			emulator.ActiveHwnd()
			shot, errShot := Screenshot(win.HWND(emulator.BindHwnd))
			shotMat, errDec := gocv.ImageToMatRGBA(shot)
			if errShot != nil || errDec != nil {
				if shotFail > 3 {
					goto Out
				}
				shotFail++
				shotMat.Close()
				TimeSleepMilli(5000)
				continue
			}

			x, y, ok := compare(shotMat, advTem)
			if ok {
				if emulator.JingjiCount < jingjiTimes {
					emulator.MouseClick(x, y-110)
					emulator.JingjiCount++
				} else {
					emulator.MouseClick(x, y)
				}

				shotMat.Close()
				TimeSleepMilli(3000)
				continue
			}

			x, y, ok = compare(shotMat, energyTem)
			if ok {
				//android.AdbShellInputTap(x, y)
				emulator.MouseClick(x, y)
				emulator.MaoxianCount--
				shotMat.Close()
				if maoxianTimes == -1 {
					TimeSleepMilli(5000)
					continue
				}
				goto Out
			}

			x, y, ok = compare(shotMat, startTem)
			if ok {
				shotMat.Close()
				if emulator.JingjiCount < jingjiTimes {
					emulator.MouseClick(40, 35)
					TimeSleepMilli(300)
					continue
				}
				if emulator.MaoxianCount < maoxianTimes || maoxianTimes == -1 {
					emulator.MouseClick(180, 260)
					TimeSleepMilli(300)
					emulator.MouseClick(x, y)
					emulator.MaoxianCount++
					TimeSleepMilli(3000)
					continue
				}
				goto Out
			}

			x, y, ok = compare(shotMat, fightTem)
			if ok {
				t := time.Now()
				re := findCard(shotMat)
				for _, v := range re {
					emulator.AdbShellInputSwipeL(v.Middle.X+30, v.Middle.Y+30, v.Middle.X+30, v.Middle.Y-130, 300)
				}

				//无牌时直接开始战斗
				if len(re) == 0 {
					emptyCardTimes++
				}
				if len(re) == 0 && emptyCardTimes >= 2 {
					//android.AdbShellInputTap(x, y)
					emulator.MouseClick(x, y)
				}

				fmt.Println("use time :", time.Since(t))
				shotMat.Close()
				TimeSleepMilli(1500)
				continue
			}
			emptyCardTimes = 0
			//android.AdbShellInputTap(940, 525)
			emulator.MouseClick(940, 525)
			shotMat.Close()
			TimeSleepMilli(2000)
			shotFail = 0
		}
	}
Out:
	emulator.Status = "结束"
	fmt.Println("end")
}

func (s *Server) getDevice(name string) *Android {
	s.RLock()
	a, ok := s.androidMap[name]
	s.RUnlock()
	if ok {
		return a
	}

	s.Lock()
	a, ok = s.androidMap[name]
	if ok {
		s.Unlock()
		return a
	}
	a = NewAndroid(name)
	s.androidMap[name] = a

	s.Unlock()

	return a
}
