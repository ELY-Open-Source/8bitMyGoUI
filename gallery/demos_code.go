package gallery

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/code"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var codeEntries = []Entry{
	{366, "Terminal", "终端", "code", "Type HELP. Up and Down bring back what you typed. XYZZY does nothing, as ever.", demoTerminal},
	{367, "CodeBlock", "代码块", "code", "Go, colored by slot: keywords, strings, numbers, builtins. COPY takes it all.", demoCodeBlock},
	{368, "LogViewer", "日志", "code", "Lines colored by level; pick the quietest worth seeing. It follows the newest.", demoLogViewer},
	{369, "HexViewer", "十六进制", "code", "A save file laid bare: offset, hex and ASCII. Arrows walk the bytes.", demoHexViewer},
	{370, "MemoryMap", "内存图", "code", "Each bank a block, filled as it fills. Red means start deleting sprites.", demoMemoryMap},
	{371, "JSONTree", "JSON树", "code", "JSON in its own order, types in their own colors. Right opens, Left closes.", demoJSONTree},
	{372, "DiffView", "差异", "code", "What changed between two versions, minus then plus.", demoDiffView},
	{373, "CrashScreen", "崩溃画面", "code", "GAME OVER in a blinking red frame, as the old guru had. CONTINUE when ready.", demoCrashScreen},
	{374, "FPSMeter", "帧率表", "code", "Frame times scroll by; green is on time, red is a stutter.", demoFPSMeter},
	{375, "ProfilerBars", "性能条", "code", "Where a frame's 16.7 ms went, nested like a family tree.", demoProfilerBars},
}

type codeDemo struct {
	shell      []string
	command    string
	least      code.Level
	logs       []code.LogLine
	byte, pick int
	node, vr   string
	crashed    bool
	marks      map[int]bool
	art        string
	pal        pixel.Palette
	frame      int
	playing    bool
	glyph      rune
	step       int
	frames     []time.Duration
	last       time.Time
	history    []string
}

func codeState(c *ui.Context) *codeDemo {
	now := pixel.K(c).Now()
	return ui.Local(c.Root(), demoKey("code"), func() codeDemo {
		return codeDemo{
			shell: []string{"8BIT-OS v1.0", "READY. TYPE HELP."}, least: code.Info,
			logs: []code.LogLine{
				{At: now.Add(-9 * time.Second), Level: code.Info, Text: "cartridge inserted"},
				{At: now.Add(-8 * time.Second), Level: code.Debug, Text: "blew on the contacts"},
				{At: now.Add(-6 * time.Second), Level: code.Warn, Text: "save battery low"},
				{At: now.Add(-3 * time.Second), Level: code.Error, Text: "princess in another castle"},
			},
			marks: map[int]bool{4: true}, art: heartArt, pal: *pixel.Current(c).Palette, glyph: 'A',
			history: []string{"load \"*\",8,1", "run", "help", "xyzzy"},
		}
	})
}

func demoTerminal(c *ui.Context, log func(string)) {
	s := codeState(c)
	r := code.Terminal(c, "8BIT-OS", s.shell, &s.command, code.TerminalOptions{Width: 260})
	if r.Changed {
		log(fmt.Sprintf("input %q", s.command))
	}
	if r.Ran != "" {
		log("ran " + r.Ran)
		s.shell = append(s.shell, "$ "+r.Ran)
		switch strings.ToLower(strings.TrimSpace(r.Ran)) {
		case "help":
			s.shell = append(s.shell, "HELP LS CAT DATE CLEAR XYZZY")
		case "ls":
			s.shell = append(s.shell, "save.dat  sprites/  music/  readme")
		case "cat readme":
			s.shell = append(s.shell, "IT'S DANGEROUS TO GO ALONE!")
		case "date":
			s.shell = append(s.shell, pixel.K(c).Now().Format("2006-01-02 15:04"))
		case "clear":
			s.shell = nil
		case "xyzzy":
			s.shell = append(s.shell, "NOTHING HAPPENS.")
		default:
			s.shell = append(s.shell, "SYNTAX ERROR")
		}
	}
}

const loop = `func main() {
	// One frame at a time, forever.
	hero := Hero{Name: "Aria", HP: 3}
	for frame := 0; ; frame++ {
		if hero.HP == 0 {
			panic("GAME OVER")
		}
		draw(hero, frame%60, nil, true)
	}
}`

func demoCodeBlock(c *ui.Context, log func(string)) {
	if code.CodeBlock(c, "main.go", loop, code.CodeBlockOptions{}).Copied {
		log("copied main.go")
	}
}

func demoLogViewer(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := codeState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		if code.LogViewer(c, "Log", s.logs, &s.least, code.LogViewerOptions{Height: 80}).Changed {
			log(fmt.Sprintf("showing level %d and up", s.least))
		}
		if input.Button(c, "+ LINE", input.ButtonOptions{Size: input.Small}).Clicked {
			lv := code.Level(len(s.logs) % 4)
			s.logs = append(s.logs, code.LogLine{At: k.Now(), Level: lv, Text: fmt.Sprintf("tick %d", len(s.logs))})
			log("added a line")
		}
	})
}

var saveFile = []byte("SAVE\x01\x00\x00\x03ARIA\x00\x00\x00\x00HP\x00\x2aGOLD\x03\xe7\x00\x00\xff\xfe\x00\x10")

func demoHexViewer(c *ui.Context, log func(string)) {
	s := codeState(c)
	if code.HexViewer(c, "save.dat", saveFile, &s.byte, code.HexViewerOptions{}).Changed {
		log(fmt.Sprintf("offset %04X", s.byte))
	}
}

func demoMemoryMap(c *ui.Context, log func(string)) {
	code.MemoryMap(c, "Cartridge memory", []code.Bank{
		{Name: "ROM0", Used: 15 << 10, Size: 16 << 10}, {Name: "ROM1", Used: 9 << 10, Size: 16 << 10},
		{Name: "WRAM", Used: 3 << 10, Size: 8 << 10}, {Name: "VRAM", Used: 6 << 10, Size: 8 << 10},
	}, code.MemoryMapOptions{})
}

const saveJSON = `{"player": {"name": "Aria", "level": 12, "alive": true,
 "items": ["sword", "potion", "map"], "pet": null}, "checkpoint": 3, "flags": {}}`

func demoJSONTree(c *ui.Context, log func(string)) {
	s := codeState(c)
	if r := code.JSONTree(c, "save.json", saveJSON, &s.node, code.JSONTreeOptions{}); r.Changed || r.Toggled {
		log("at " + s.node)
	}
}

func demoDiffView(c *ui.Context, log func(string)) {
	before := "speed = 4\njump = 2\nlives = 3\ncoins = 0"
	after := "speed = 4\njump = 3\nlives = 3\ncoins = 0\nwarp = true"
	code.DiffView(c, "settings.cfg", before, after, code.DiffViewOptions{})
}

func demoCrashScreen(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := codeState(c)
	if s.crashed {
		if input.Button(c, "CRASH AGAIN", input.ButtonOptions{}).Clicked {
			s.crashed = false
			log("crashed again")
		}
		return
	}
	stack := []string{"main.loot(0x3)", "main.openChest(...)", "main.main()"}
	if code.CrashScreen(c, "panic: index out of range [3] with length 3", stack, code.CrashScreenOptions{Width: 280}).Continue {
		s.crashed = true
		log("continued at " + k.Now().Format("15:04:05"))
	}
}

func demoFPSMeter(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := codeState(c)
	// Made-up frames: mostly on time, now and then a hitch.
	add := func(i int) {
		ms := 9 + 4*math.Sin(float64(i)/5)
		if i%23 == 7 {
			ms = 38
		}
		s.frames = append(s.frames[max(0, len(s.frames)-200):], time.Duration(ms*float64(time.Millisecond)))
	}
	if len(s.frames) == 0 {
		for i := range 56 {
			add(i)
		}
	}
	if !k.Reduce() {
		if t := k.Tick(c, 12); s.last.IsZero() || k.Now().Sub(s.last) >= time.Second/12 {
			s.last = k.Now()
			add(t)
		}
	}
	code.FPSMeter(c, "Frame times", s.frames, code.FPSMeterOptions{})
}

func demoProfilerBars(c *ui.Context, log func(string)) {
	ms := func(f float64) time.Duration { return time.Duration(f * float64(time.Millisecond)) }
	code.ProfilerBars(c, "One frame", code.Span{Name: "frame", Length: ms(16.7), Children: []code.Span{
		{Name: "input", Length: ms(1)},
		{Name: "update", Start: ms(1), Length: ms(6), Children: []code.Span{
			{Name: "physics", Start: ms(1), Length: ms(3)}, {Name: "ai", Start: ms(4), Length: ms(2.5)},
		}},
		{Name: "render", Start: ms(7.5), Length: ms(8), Children: []code.Span{
			{Name: "sprites", Start: ms(7.5), Length: ms(5), Children: []code.Span{{Name: "dither", Start: ms(9), Length: ms(2)}}},
			{Name: "text", Start: ms(12.5), Length: ms(2)},
		}},
	}}, code.ProfilerBarsOptions{Width: 260})
}
