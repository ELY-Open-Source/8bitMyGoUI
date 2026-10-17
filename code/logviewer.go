package code

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Level is how loud a log line is.
type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

var levels = []struct {
	name string
	slot pixel.Slot
}{{"DEBUG", pixel.Faint}, {"INFO", pixel.Ink}, {"WARN", pixel.Warn}, {"ERROR", pixel.Bad}}

type LogLine struct {
	At    time.Time
	Level Level
	Text  string
}

type LogViewerOptions struct {
	// Width and Height are dots; zeros mean 260 by 120.
	Width, Height int
	Disabled      bool
}

type LogViewerResult struct {
	Element ui.Element
	// Changed is the level filter.
	Changed bool
}

type logKey struct{}

// LogViewer colors lines by level, hides those under *least, follows.
func LogViewer(c *ui.Context, text string, lines []LogLine, least *Level, o LogViewerOptions) LogViewerResult {
	k := pixel.K(c)
	if *least < Debug || *least > Error {
		panic(fmt.Sprintf("code: LogViewer level %d", *least))
	}
	for _, l := range lines {
		if l.Level < Debug || l.Level > Error {
			panic(fmt.Sprintf("code: LogViewer line %q has level %d", l.Text, l.Level))
		}
	}
	w, h := pixel.Or("LogViewer Width", o.Width, 260), pixel.Or("LogViewer Height", o.Height, 120)
	names := make([]string, len(levels))
	for i, l := range levels {
		names[i] = l.name
	}
	var res LogViewerResult
	res.Element = ui.Column(c).Width(k.D(w)).Gap(k.D(3)).Children(func() {
		res.Changed = input.Segmented(c, text, names, (*int)(least), input.SegmentedOptions{Disabled: o.Disabled}).Changed
		layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well, Padding: 2}, func() {
			box := ui.Box(c).FillWidth().Height(k.D(h))
			s := ui.Local(box, logKey{}, func() ui.ScrollState { return ui.ScrollState{} })
			box.Children(func() {
				layout.ScrollArea(c, s, layout.ScrollAreaOptions{Label: text, Follow: true, Disabled: o.Disabled}, func() {
					ui.Column(c).FillWidth().Children(func() {
						for _, l := range lines {
							if l.Level >= *least {
								row(c, k, l)
							}
						}
					})
				}).Element.Fill()
			})
		})
	})
	if res.Changed {
		slog.Debug("8bit: log viewer", "label", text, "least", levels[*least].name)
	}
	return res
}

// row is a line: time, level, words in its slot.
func row(c *ui.Context, k pixel.Kit, l LogLine) {
	lv, at := levels[l.Level], l.At.Format("15:04:05")
	ui.Row(c).FillWidth().Gap(k.D(4)).Role(ui.RoleText).Label(at + " " + lv.name + " " + l.Text).Children(func() {
		ui.Box(c).Size(k.D(56), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Micro(0, 4, at, pixel.Faint, 1)
			g.Micro(34, 4, lv.name, lv.slot, 1)
		})
		pixel.Text(c, l.Text, pixel.TextOptions{Slot: lv.slot, Wrap: true}).Grow(1)
	})
}
