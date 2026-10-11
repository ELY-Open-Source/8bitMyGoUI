// Package gallery shows every component live, in every palette.
package gallery

import (
	"embed"
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Entry is one page: a component or an engine foundation.
type Entry struct {
	// ID is 1–500 for components, 0 for foundations.
	ID       int
	Name, Zh string
	Category string
	Summary  string
	Demo     func(c *ui.Context, log func(string))
}

type Category struct {
	Key, Name, Zh string
	Icon          *pixel.Sprite
}

var Categories = []Category{
	{"foundations", "Foundations", "基础", icon.Palette},
	{"layout", "Layout", "布局", icon.Grid},
	{"text", "Text", "文字", icon.FileText},
	{"display", "Display", "展示", icon.Image},
	{"input", "Input", "输入", icon.Edit},
	{"feedback", "Feedback", "反馈", icon.Info},
	{"navigation", "Navigation", "导航", icon.Map},
	{"data", "Data", "数据", icon.List},
	{"overlay", "Overlay", "浮层", icon.Copy},
	{"datetime", "Date & Time", "日期时间", icon.Calendar},
	{"chart", "Chart", "图表", icon.Sort},
	{"game", "Game", "游戏", icon.Gamepad},
	{"fx", "Effects", "特效", icon.Sparkle},
	{"chat", "Chat", "聊天", icon.Chat},
	{"media", "Media", "媒体", icon.Music},
	{"code", "Code", "代码", icon.Code},
	{"account", "Account", "账户", icon.User},
	{"files", "Files", "文件", icon.Folder},
	{"project", "Project", "项目", icon.Flag},
	{"shop", "Shop", "商店", icon.Cart},
	{"agent", "Agent", "智能体", icon.Robot},
	{"desk", "Desk", "桌面", icon.Save},
	{"tabletop", "Tabletop", "桌游", icon.Crown},
}

var entries = slices.Concat(foundationEntries, layoutEntries, textEntries)

func Entries() []Entry { return entries }

func (e Entry) Slug() string {
	s := strings.ToLower(strings.ReplaceAll(e.Name, " ", "-"))
	if e.ID == 0 {
		return s
	}
	return fmt.Sprintf("%03d-%s", e.ID, s)
}

func (e Entry) Path() string {
	if e.ID == 0 {
		return "/f/" + e.Slug()
	}
	return "/c/" + e.Slug()
}

//go:embed demos_*.go
var demoFiles embed.FS

var codes = sources()

// Code is the source of the entry's demo function.
func (e Entry) Code() string {
	src, ok := codes[funcName(e.Demo)]
	if !ok {
		panic("gallery: no source for " + funcName(e.Demo))
	}
	return src
}

func funcName(f any) string {
	n := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return n[strings.LastIndex(n, ".")+1:]
}

// sources indexes every top-level demo function by name.
func sources() map[string]string {
	files, err := demoFiles.ReadDir(".")
	if err != nil {
		panic(err)
	}
	out := map[string]string{}
	for _, f := range files {
		b, err := demoFiles.ReadFile(f.Name())
		if err != nil {
			panic(err)
		}
		for _, part := range strings.Split(string(b), "\nfunc ")[1:] {
			end := strings.Index(part, "\n}\n")
			if end < 0 {
				panic("gallery: unterminated func in " + f.Name())
			}
			if name := part[:strings.IndexAny(part, "( ")]; name != "" {
				out[name] = "func " + part[:end+2]
			}
		}
	}
	return out
}
