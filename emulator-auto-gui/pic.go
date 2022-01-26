package main

import (
	"errors"
	"fmt"
	"github.com/lxn/win"
	"github.com/vcaesar/gcv"
	"gocv.io/x/gocv"
	"image"
	"reflect"
	"unsafe"
)

var CardTemplate []gocv.Mat

func init() {
	CardTemplate = cardTemplate()
}

func findCard(imgSrc gocv.Mat) map[gcv.Point]gcv.Result {
	m := make(map[gcv.Point]gcv.Result)
	for _, v := range CardTemplate {
		resultSlice := gcv.FindAllTemplate(imgSrc, v, 0.8)

		for _, j := range resultSlice {
			flag := 0
			for mk, _ := range m {
				if ((j.Middle.X + 10) > mk.X) && ((j.Middle.X - 10) < mk.X) && ((j.Middle.Y + 10) > mk.Y) && ((j.Middle.Y - 10) < mk.Y) {
					flag = 1
				}
			}
			if flag == 1 {
				continue
			}
			m[j.Middle] = j
		}

	}
	return m
}

func compare(imgSrc, imgSearch gocv.Mat) (x, y int, success bool) {
	resultSlice := gcv.FindAllTemplate(imgSrc, imgSearch, 0.9)
	var tmp float32
	for _, v := range resultSlice {
		if v.MaxVal[0] > tmp {
			tmp = v.MaxVal[0]
			x = v.Middle.X
			y = v.Middle.Y
			success = true
		}
	}

	fmt.Printf("x : %d,y : %d, maxval : %f \n", x, y, tmp)

	return
}

func cardTemplate() (cardTem []gocv.Mat) {
	for i := 1; i < 7; i++ {
		name := fmt.Sprintf("static/card/one%d.png", i)
		cardTem = append(cardTem, loadTemplate(name))
	}
	for i := 1; i < 5; i++ {
		name := fmt.Sprintf("static/card/zero%d.png", i)
		cardTem = append(cardTem, loadTemplate(name))
	}

	return
}

func loadTemplate(picName string) gocv.Mat {
	pic, err := fs.ReadFile(picName)
	if err != nil {
		panic(err)
	}
	picMat, err := gocv.IMDecode(pic, gocv.IMReadFlag(4))
	if err != nil {
		panic(err)
	}
	return picMat
}

func Screenshot(hwnd win.HWND) (*image.RGBA, error) {
	hdc := win.GetDC(hwnd)
	if hdc == 0 {
		return nil, errors.New("get hdc error")
	}
	defer win.ReleaseDC(hwnd, hdc)

	memHDC := win.CreateCompatibleDC(hdc)
	if memHDC == 0 {
		return nil, errors.New("create compatible dc err")
	}
	defer win.DeleteDC(memHDC)

	var r win.RECT
	if !win.GetWindowRect(hwnd, &r) {
		return nil, errors.New("GetWindowRect failed")
	}

	width, height := r.Right-r.Left, r.Bottom-r.Top

	var bt win.BITMAPINFO
	bt.BmiHeader.BiSize = uint32(unsafe.Sizeof(bt.BmiHeader))
	bt.BmiHeader.BiWidth = width
	bt.BmiHeader.BiHeight = -height
	bt.BmiHeader.BiPlanes = 1
	bt.BmiHeader.BiBitCount = 32
	bt.BmiHeader.BiCompression = win.BI_RGB

	var ptr unsafe.Pointer
	memHBMP := win.CreateDIBSection(memHDC, &bt.BmiHeader, win.DIB_RGB_COLORS, &ptr, 0, 0)
	if memHBMP == 0 {
		return nil, errors.New("create DIB section err")
	}
	if memHBMP == 2 {
		return nil, errors.New("one or more of the input parameters is invalid while calling CreateDIBSection")
	}
	defer win.DeleteObject(win.HGDIOBJ(memHBMP))

	obj := win.SelectObject(memHDC, win.HGDIOBJ(memHBMP))
	if obj == 0 {
		return nil, errors.New("error occurred and the selected object is not a region")
	}
	if obj == 0xffffffff { //GDI_ERROR
		return nil, errors.New("GDI_ERROR while calling SelectObject")
	}
	defer win.DeleteObject(obj)

	if !win.BitBlt(memHDC, 0, 0, width, height, hdc, 0, 0, win.SRCCOPY) {
		return nil, errors.New("BitBlt failed err")
	}

	var slice []byte
	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	hdrp.Data = uintptr(ptr)
	hdrp.Len = int(width * height * 4)
	hdrp.Cap = int(width * height * 4)

	imageBytes := make([]byte, len(slice))
	for i := 0; i < len(imageBytes); i += 4 {
		imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
	}
	img := &image.RGBA{
		Pix:    imageBytes,
		Stride: int(4 * width),
		Rect:   image.Rect(0, 0, int(width), int(height)),
	}
	return img, nil
}
