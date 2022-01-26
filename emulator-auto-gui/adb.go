package main

import (
	"bytes"
	"errors"
	"github.com/lxn/win"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	APPPackageName = "com.game.game"
	ClassName      = "com.game.start"
)

var AdbPath = `D:\BaiduNetdiskDownload\QtScrcpy-win-x64-v1.7.1/adb.exe`
var EmulatorPath = ``

//var AdbPath = "./adb.exe"

type Android struct {
	Name string
}

func NewAndroid(name string) *Android {
	a := &Android{name}
	return a
}

func newCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

func (a *Android) AdbShellInputTap(x, y int) {
	x2 := strconv.Itoa(x)
	y2 := strconv.Itoa(y)
	newCmd(AdbPath, "-s", a.Name, "shell", "input", "tap", x2, y2).Run()
	//exec.Command(AdbPath, "-s", a.Name, "shell", "input", "tap", x2, y2).Run()
}

func (a *Android) AdbShellAmStart(program string) {
	newCmd(AdbPath, "-s", a.Name, "shell", "am", "start", program).Run()
	//exec.Command(AdbPath, "-s", a.Name, "shell", "am", "start", program).Run()
}

//模拟滑动
//adb shell input swipe  0 0  600 600
func (a *Android) AdbShellInputSwipe(x1, y1, x2, y2 int) {
	xx1 := strconv.Itoa(x1)
	yy1 := strconv.Itoa(y1)
	xx2 := strconv.Itoa(x2)
	yy2 := strconv.Itoa(y2)
	newCmd(AdbPath, "-s", a.Name, "shell", "input", "swipe", xx1, yy1, xx2, yy2).Run()
	//exec.Command(AdbPath, "-s", a.Name, "shell", "input", "swipe", xx1, yy1, xx2, yy2).Run()
}

//模拟长按 最后一个参数1000表示1秒，可将下面某个参数由500改为501，即允许坐标点有很小的变化。
//adb shell input swipe  500 500  500 500 1000
func (a *Android) AdbShellInputSwipeL(x1, y1, x2, y2, t int) {
	xx1 := strconv.Itoa(x1)
	yy1 := strconv.Itoa(y1)
	xx2 := strconv.Itoa(x2)
	yy2 := strconv.Itoa(y2)
	t1 := strconv.Itoa(t)
	newCmd(AdbPath, "-s", a.Name, "shell", "input", "swipe", xx1, yy1, xx2, yy2, t1).Run()
	//exec.Command(AdbPath, "-s", a.Name, "shell", "input", "swipe", xx1, yy1, xx2, yy2, t1).Run()
	//fmt.Println("swipe err: ",err)
}

func (a *Android) PullScreenShot() (img []byte, err error) {
	cmd := newCmd(AdbPath, "-s", a.Name, "shell", "screencap", "-p")
	//cmd := exec.Command(AdbPath, "-s", a.Name, "shell", "screencap", "-p")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err = cmd.Run(); err != nil {
		return nil, err
	}
	x := bytes.Replace(out.Bytes(), []byte("\r\n"), []byte("\n"), -1)
	//img, err = png.Decode(bytes.NewReader(x))
	return x, nil

}

func GetAllDevices() []string {
	cmd := newCmd(AdbPath, "devices")
	//cmd := exec.Command(AdbPath, "devices")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	list := strings.Split(out.String(), "\n")

	var devices []string
	for k, _ := range list {

		if strings.Contains(list[k], "emulator") {
			i := strings.Index(list[k], "\t")
			if i > -1 {
				devices = append(devices, list[k][:i])
			}
		}

	}
	return devices
}

//等待x毫秒
func TimeSleepMilli(x int) {
	time.Sleep(time.Duration(x) * time.Millisecond)
}

type Emulator struct {
	Title    string
	Index    int
	TopHwnd  int
	BindHwnd int

	Status       string
	JingjiCount  int
	MaoxianCount int
	checked      bool
}

func GetEmulatorList() ([]*Emulator, error) {
	cmd := newCmd(path.Join(EmulatorPath, "ldconsole.exe"), "list2")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	str := out.String()
	if len(str) < 1 {
		return nil, errors.New("out err")
	}

	uStr, err := simplifiedchinese.GBK.NewDecoder().String(str)
	if err != nil {
		return nil, err
	}

	var emulatorList []*Emulator
	list := strings.Split(uStr, "\r\n")
	for i := range list {
		tmp := strings.Split(list[i], ",")
		if len(tmp) == 7 && tmp[4] == "1" {
			emulator := Emulator{Title: tmp[1]}
			emulator.Index, _ = strconv.Atoi(tmp[0])
			emulator.TopHwnd, _ = strconv.Atoi(tmp[2])
			emulator.BindHwnd, _ = strconv.Atoi(tmp[3])
			emulatorList = append(emulatorList, &emulator)
		}
	}
	return emulatorList, nil
}

func (e *Emulator) StartApp(program string) {
	index := strconv.Itoa(e.Index)
	newCmd(path.Join(EmulatorPath, "ld.exe"), "-s", index, "am", "start", program).Run()
}

func (e *Emulator) AdbShellInputSwipeL(x1, y1, x2, y2, t int) {
	xx1 := strconv.Itoa(x1)
	yy1 := strconv.Itoa(y1)
	xx2 := strconv.Itoa(x2)
	yy2 := strconv.Itoa(y2)
	t1 := strconv.Itoa(t)
	index := strconv.Itoa(e.Index)
	newCmd(path.Join(EmulatorPath, "ld.exe"), "-s", index, "input", "swipe", xx1, yy1, xx2, yy2, t1).Run()
	//exec.Command(AdbPath, "-s", a.Name, "shell", "input", "swipe", xx1, yy1, xx2, yy2, t1).Run()
	//fmt.Println("swipe err: ",err)
}

func (e *Emulator) ActiveHwnd() {
	topHwnd := win.HWND(e.TopHwnd)
	if win.IsIconic(topHwnd) {
		win.ShowWindow(topHwnd, win.SW_RESTORE)
	}
	rect := win.RECT{}
	if win.GetWindowRect(topHwnd, &rect) {
		if (rect.Right-rect.Left) != 1002 || (rect.Bottom-rect.Top) != 576 {
			win.MoveWindow(topHwnd, 100, 100, 1002, 576, true)
		}
	}
}

func (e *Emulator) MouseClick(x, y int) {
	win.PostMessage(win.HWND(e.BindHwnd), win.WM_LBUTTONDOWN, 0, uintptr(x+(y<<16)))
	win.PostMessage(win.HWND(e.BindHwnd), win.WM_LBUTTONUP, 0, 0)
}

func (e *Emulator) MouseDrag(x1, y1, x2, y2, delay int) {
	/*win.SendMessage(win.HWND(e.BindHwnd), win.WM_LBUTTONDOWN, win.MK_LBUTTON, uintptr(x1+(y1<<16)))
	win.SendMessage(win.HWND(e.BindHwnd), win.WM_LBUTTONUP, win.MK_LBUTTON, uintptr(x1+(y1<<16)))
	win.SendMessage(win.HWND(e.BindHwnd), win.WM_LBUTTONDOWN, win.MK_LBUTTON, uintptr(x1+(y1<<16)))
	TimeSleepMilli(delay+200)*/
	/*for i := 1; i <= 10; i++ {
		pos := uintptr(x1 + (x2-x1)/10 + ((y1 + (y2-y1)/10) << 16))
		win.PostMessage(win.HWND(e.BindHwnd), win.WM_MOUSEMOVE, win.MK_LBUTTON, pos)
		TimeSleepMilli(delay / 10)
	}*/

	/*win.SendMessage(win.HWND(e.BindHwnd), win.WM_MOUSEMOVE, win.MK_LBUTTON, uintptr(x2+(y2<<16)))
	TimeSleepMilli(delay)
	win.SendMessage(win.HWND(e.BindHwnd), win.WM_LBUTTONUP, 0, uintptr(x2+(y2<<16)))*/

	win.SendMessage(win.HWND(e.BindHwnd), win.WM_LBUTTONDOWN, 0, uintptr(x1+(y1<<16)))
	TimeSleepMilli(delay + 700)
	win.SendMessage(win.HWND(e.BindHwnd), win.WM_MOUSEMOVE, 0, uintptr(x2+(y2<<16)))
	TimeSleepMilli(delay)
	win.SendMessage(win.HWND(e.BindHwnd), win.WM_LBUTTONUP, 0, 0)
}
