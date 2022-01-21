package main

import (
	"embed"
	"fmt"
	"github.com/lxn/win"
	"gocv.io/x/gocv"
	"log"
	"sync"
	"time"
)

//go:embed static
var fs embed.FS

var server = &Server{androidMap: make(map[string]*Android)}

var DeviceList []string

type Server struct {
	sync.RWMutex
	sync.WaitGroup
	androidMap   map[string]*Android
	emulatorList []Emulator
}

func main() {
	/*test :=Emulator{BindHwnd: 2034064}
	test.MouseDrag(31,288,310,288,1000)
	return*/
	go func() {
		list, err := GetEmulatorList()
		if err != nil {
			log.Fatalln("get emulator list err:", err)
		}
		server.emulatorList = list
		fmt.Printf("%#v\n", server)
		for k := range server.emulatorList {
			//server.androidMap[k].AdbShellInputSwipeL(120,460,120,260,5000)
			//return
			server.Add(1)
			go server.start(server.emulatorList[k])
		}
		server.Wait()

	}()

	out := gocv.NewWindow("StarSharks")
	defer out.Close()
	out.WaitKey(0)

	fmt.Println("finish")
}

func (s *Server) start(emulator Emulator) {
	defer s.Done()

	emulator.StartApp(APPPackageName + "/" + ClassName)

	advTem := loadTemplate("static/adventure.png")
	startTem := loadTemplate("static/start.png")
	energyTem := loadTemplate("static/energyless.png")
	fightTem := loadTemplate("static/fight.png")

	emptyCardTimes := 0
	shotFail := 0
	for {
		emulator.ActiveHwnd()
		shot, errShot := Screenshot(win.HWND(emulator.BindHwnd))
		shotMat, errDec := gocv.ImageToMatRGBA(shot)
		if errShot != nil || errDec != nil {
			if shotFail > 3 {
				break
			}
			shotFail++
			shotMat.Close()
			TimeSleepMilli(5000)
			continue
		}

		x, y, ok := compare(shotMat, advTem)
		if ok {
			//android.AdbShellInputTap(x, y)
			emulator.MouseClick(x, y)
			shotMat.Close()
			TimeSleepMilli(3000)
			continue
		}

		x, y, ok = compare(shotMat, energyTem)
		if ok {
			//android.AdbShellInputTap(x, y)
			emulator.MouseClick(x, y)
			shotMat.Close()
			continue
		}

		x, y, ok = compare(shotMat, startTem)
		if ok {
			//android.AdbShellInputTap(180, 260)
			emulator.MouseClick(180, 260)
			TimeSleepMilli(300)
			//android.AdbShellInputTap(x, y)
			emulator.MouseClick(x, y)
			shotMat.Close()
			TimeSleepMilli(3000)
			continue
		}

		x, y, ok = compare(shotMat, fightTem)
		if ok {
			/*var startx, starty, i = 120, 460, 0
			t := time.Now()
			for i < 11 {
				android.AdbShellInputSwipeL(startx+63*i, starty, startx+63*i, starty-100, 1)
				TimeSleepMilli(100)
				i++
			}
			fmt.Println("use time :", time.Since(t))
			android.AdbShellInputTap(x, y)
			fmt.Println("tap flight end")*/
			t := time.Now()
			re := findCard(shotMat)
			for _, v := range re {
				emulator.AdbShellInputSwipeL(v.Middle.X+30, v.Middle.Y+30, v.Middle.X+30, v.Middle.Y-130, 300)
				//emulator.MouseDrag(v.Middle.X+30, v.Middle.Y+45, v.Middle.X+30, v.Middle.Y-330, 500)
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
