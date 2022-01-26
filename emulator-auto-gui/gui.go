package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"math/rand"
	"os"
	"time"
)

func (s *Server) newGui() {
	rand.Seed(time.Now().UnixNano())
	mw := s.mW
	model := new(StatusModel)
	model.items = server.emulatorList

	if _, err := (MainWindow{
		Icon:     "my.ico",
		AssignTo: &mw.MainWindow,
		Title:    "StarSharks-Auto",

		Size:   Size{Width: 500, Height: 850},
		Layout: VBox{MarginsZero: true},

		Children: []Widget{
			Composite{
				MaxSize: Size{Height: 50},
				Layout:  HBox{Spacing: 10, Margins: Margins{Left: 10, Top: 10, Right: 250, Bottom: 10}},
				Children: []Widget{
					Label{
						Text:    "雷电模拟器路径:",
						MaxSize: Size{Width: 90},
						MinSize: Size{Width: 90},
						//Font:    Font{PointSize: 10, Bold: true},
					},
					//HSpacer{},
					LineEdit{
						AssignTo:    &mw.EmulatorPath,
						MinSize:     Size{Width: 300},
						MaxSize:     Size{Width: 300},
						Text:        EmulatorPath,
						TextColor:   walk.RGB(255, 0, 0),
						ToolTipText: "输入模拟器正确路径",
					},
					//HSpacer{},
					PushButton{
						//MaxSize: Size{Width: 50,Height: 1},
						Text:     "刷新模拟器列表",
						AssignTo: &mw.RefreshList,
						MaxSize:  Size{Width: 100, Height: 100},
						OnClicked: func() {
							go func() {
								EmulatorPath = mw.EmulatorPath.Text()
								list, err := GetEmulatorList()
								if err == nil {
									os.WriteFile("emulator-path.txt", []byte(EmulatorPath), 0600)
								}
								server.emulatorList = list
								model.items = list
								mw.InfoTable.SetModel(model)
								model.PublishRowsReset()
							}()
						},
					},
					HSpacer{},
				},
			},
			VSpacer{},
			Composite{
				MaxSize: Size{Height: 50},
				Layout:  HBox{Spacing: 10, Margins: Margins{Left: 10, Top: 10, Right: 250, Bottom: 10}},
				Children: []Widget{
					Label{
						Text:    "竞技模式:",
						MaxSize: Size{Width: 90},
						MinSize: Size{Width: 90},
						//Font:    Font{PointSize: 10, Bold: true},
					},
					//HSpacer{},
					CheckBox{
						AssignTo: &mw.JingjiCheck,
						Text:     "刷",
						MaxSize:  Size{Width: 40},
						MinSize:  Size{Width: 40},
					},
					//HSpacer{},
					NumberEdit{
						AssignTo:  &mw.JingjiNum,
						MinSize:   Size{Width: 60},
						MaxSize:   Size{Width: 60},
						Suffix:    "次",
						Value:     Bind("2"),
						TextColor: walk.RGB(255, 0, 0),
					},
					HSpacer{},
				},
			},
			VSpacer{},
			Composite{
				MaxSize: Size{Height: 50},
				Layout:  HBox{Spacing: 10, Margins: Margins{Left: 10, Top: 10, Right: 110, Bottom: 10}},
				Children: []Widget{
					Label{
						Text:    "冒险模式:",
						MaxSize: Size{Width: 90},
						MinSize: Size{Width: 90},
						//Font:    Font{PointSize: 10, Bold: true},
					},
					//HSpacer{},
					CheckBox{
						AssignTo: &mw.MaoXianCheck,
						Text:     "刷",
						MaxSize:  Size{Width: 40},
						MinSize:  Size{Width: 40},
					},
					//HSpacer{},
					NumberEdit{
						AssignTo:  &mw.MaoXianNum,
						MinSize:   Size{Width: 60},
						MaxSize:   Size{Width: 60},
						Suffix:    "次",
						Value:     Bind("8"),
						TextColor: walk.RGB(255, 0, 0),
					},
					//HSpacer{},
					CheckBox{
						AssignTo: &mw.MaoXianEnd,
						Text:     "冒险至没活力",
						MaxSize:  Size{Width: 100},
					},
					HSpacer{},
				},
			},

			Composite{
				MaxSize: Size{Height: 50},
				Layout:  HBox{Spacing: 10, Margins: Margins{Left: 10, Top: 10, Right: 110, Bottom: 10}},
				Children: []Widget{
					CheckBox{
						AssignTo: &mw.AllCheck,
						Text:     "全选",
						OnClicked: func() {
							mw.AllCheck.SetEnabled(false)
							defer mw.AllCheck.SetEnabled(true)
							for _, v := range model.items {
								v.checked = mw.AllCheck.Checked()
							}
							mw.InfoTable.SetModel(model)
							//model.PublishRowsReset()
						},
					},
					//HSpacer{},
					PushButton{
						//MaxSize: Size{Width: 50,Height: 1},
						Text:     "开始任务",
						AssignTo: &mw.StartTask,
						MaxSize:  Size{Width: 100, Height: 100},
						OnClicked: func() {
							go func() {
								mw.StartTask.SetEnabled(false)
								defer mw.StartTask.SetEnabled(true)
								if mw.StartTask.Text() == "开始任务" {
									mw.AllCheck.SetEnabled(false)
									mw.RefreshList.SetEnabled(false)
									mw.InfoTable.SetCheckBoxes(false)
									mw.StartTask.SetText("结束任务")
									go s.startTask()
								} else {
									mw.AllCheck.SetEnabled(true)
									mw.RefreshList.SetEnabled(true)
									mw.InfoTable.SetCheckBoxes(true)
									mw.StartTask.SetText("开始任务")
									mw.InfoTable.SetModel(model)
									close(server.closeCh)
									s.Wait()
								}
							}()
						},
					},
					HSpacer{},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true},
				//Layout: Grid{Columns: 5, Spacing: 10},
				Children: []Widget{
					TableView{
						AssignTo:         &mw.InfoTable,
						AlternatingRowBG: true,
						CheckBoxes:       true,
						ColumnsOrderable: true,
						MultiSelection:   true,
						Columns: []TableViewColumn{
							{Title: "序号", Frozen: true, Alignment: AlignCenter, Width: 80},
							{Title: "模拟器名称", Frozen: true, Width: 120},
							{Title: "当前状态", Frozen: true, Width: 120},
							{Title: "竞技", Frozen: true, Width: 80},
							{Title: "冒险", Frozen: true, Width: 80},
						},
						Model: model,
						StyleCell: func(style *walk.CellStyle) {
							item := model.items[style.Row()]

							if item.checked {
								if style.Row()%2 == 0 {
									style.BackgroundColor = walk.RGB(255, 215, 255)
								} else {
									style.BackgroundColor = walk.RGB(220, 199, 239)
								}
							}
						},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatalln(err)
	}
}

type MyMainWindow struct {
	*walk.MainWindow
	RefreshList  *walk.PushButton
	EmulatorPath *walk.LineEdit
	StartTask    *walk.PushButton
	JingjiCheck  *walk.CheckBox
	MaoXianCheck *walk.CheckBox
	JingjiNum    *walk.NumberEdit
	MaoXianNum   *walk.NumberEdit
	MaoXianEnd   *walk.CheckBox
	AllCheck     *walk.CheckBox
	InfoTable    *walk.TableView
	CurrentTitle string
}

type Info struct {
	Index         int
	VmName        string
	Status        string
	JingJiStatus  string
	MaoXianStatus string
	checked       bool
}

type StatusModel struct {
	walk.TableModelBase
	items []*Emulator
}

func (t *StatusModel) RowCount() int {
	return len(t.items)
}
func (t *StatusModel) Value(row, col int) interface{} {
	item := t.items[row]

	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Title
	case 2:
		return item.Status
	case 3:
		return fmt.Sprintf("%d次", item.JingjiCount)
	case 4:
		return fmt.Sprintf("%d次", item.MaoxianCount)
	}
	panic("unexpected col")
}
func (t *StatusModel) Checked(row int) bool {
	return t.items[row].checked
}
func (t *StatusModel) SetChecked(row int, checked bool) error {
	t.items[row].checked = checked
	return nil
}
