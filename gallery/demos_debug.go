package gallery

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/code"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var debugEntries = []Entry{
	{376, "Breakpoints", "断点", "code", "Click a line, or press Enter on it, to stop there next time.", demoBreakpoints},
	{377, "Variables", "变量", "code", "Names, types and values; open a struct to see inside.", demoVariables},
	{378, "CommandHistory", "命令历史", "code", "Everything you typed, numbered. Enter runs it again.", demoCommandHistory},
	{379, "ShortcutList", "快捷键表", "code", "Actions beside their keycaps. Do not trust the last one.", demoShortcutList},
	{380, "DebugGrid", "调试网格", "code", "The dot grid laid over anything, with rulers in dots.", demoDebugGrid},
	{381, "PixelEditor", "像素编辑器", "code", "Pen, fill, eraser, picker. Arrows and Enter paint too. The art below is live.", demoPixelEditor},
	{382, "PaletteEditor", "调色板编辑", "code", "Sixteen colors, sixteen levels a channel, like the old chips. The preview keeps up.", demoPaletteEditor},
	{383, "SpriteSheet", "精灵表", "code", "Every frame in a grid; PLAY spins the coin.", demoSpriteSheet},
	{384, "GlyphMap", "字形表", "code", "128 code points a page; PageUp and PageDown turn it. Dots mark what the font lacks.", demoGlyphMap},
	{385, "RegisterView", "寄存器", "code", "A 6502's registers and flags. STEP runs a line; what it changed lights up.", demoRegisterView},
}

const routine = `func jump(h *Hero) {
	if h.OnGround {
		h.VY = -4
	}
	h.Y += h.VY
}`

func demoBreakpoints(c *ui.Context, log func(string)) {
	s := codeState(c)
	if code.Breakpoints(c, "jump.go", routine, s.marks, code.BreakpointsOptions{}).Changed {
		log(fmt.Sprintf("breakpoints on %d lines", len(s.marks)))
	}
}

func demoVariables(c *ui.Context, log func(string)) {
	s := codeState(c)
	vars := []code.Var{
		{Name: "hero", Type: "Hero", Value: "{…}", Fields: []code.Var{
			{Name: "Name", Type: "string", Value: `"Aria"`}, {Name: "HP", Type: "int", Value: "42"},
			{Name: "Pos", Type: "Vec", Value: "{…}", Fields: []code.Var{{Name: "X", Type: "float64", Value: "12.5"}, {Name: "Y", Type: "float64", Value: "-3"}}},
			{Name: "Flying", Type: "bool", Value: "false"},
		}},
		{Name: "target", Type: "*Enemy", Value: "nil"},
		{Name: "lives", Type: "uint8", Value: "3"},
	}
	if r := code.Variables(c, "Locals", vars, &s.vr, code.VariablesOptions{Width: 220}); r.Changed || r.Toggled {
		log("at " + s.vr)
	}
}

func demoCommandHistory(c *ui.Context, log func(string)) {
	s := codeState(c)
	r := code.CommandHistory(c, "History", s.history, &s.pick, code.CommandHistoryOptions{Height: 70})
	if r.Rerun {
		log("rerun " + s.history[s.pick])
	} else if r.Changed {
		log("picked " + s.history[s.pick])
	}
}

func demoShortcutList(c *ui.Context, log func(string)) {
	code.ShortcutList(c, "Controls", []code.Shortcut{
		{Action: "Save", Keys: []string{"Ctrl", "S"}}, {Action: "Jump", Keys: []string{"Space"}},
		{Action: "Pause", Keys: []string{"Esc"}},
		{Action: "Free gold", Keys: []string{"Alt", "F4"}},
	}, code.ShortcutListOptions{Width: 300})
}

func demoDebugGrid(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	code.DebugGrid(c, code.DebugGridOptions{}, func() {
		ui.Column(c).Padding(k.D(8)).Gap(k.D(6)).Children(func() {
			if input.Button(c, "ALIGNED?", input.ButtonOptions{}).Clicked {
				log("pressed on the grid")
			}
			pixel.Text(c, "Every edge on a dot.", pixel.TextOptions{})
		})
	})
}

const heartArt = `..rr..rr..
.rrrrrrrr.
rrlrrrrrrr
rrrrrrrrrr
.rrrrrrrr.
..rrrrrr..
...rrrr...
....rr....
`

func demoPixelEditor(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := codeState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		r := code.PixelEditor(c, "Heart", &s.art, code.PixelEditorOptions{})
		if r.Changed {
			log("painted")
		}
		if r.Copied {
			log("art copied")
		}
		for _, l := range strings.Split(strings.TrimSpace(s.art), "\n") {
			pixel.Text(c, l, pixel.TextOptions{Slot: pixel.Faint})
		}
	})
}

func demoPaletteEditor(c *ui.Context, log func(string)) {
	s := codeState(c)
	if code.PaletteEditor(c, "Palette copy", &s.pal, code.PaletteEditorOptions{}).Changed {
		log("palette edited")
	}
}

func demoSpriteSheet(c *ui.Context, log func(string)) {
	s := codeState(c)
	if code.SpriteSheet(c, "Coin", coin, &s.frame, &s.playing, code.SpriteSheetOptions{Scale: 3}).Changed {
		log(fmt.Sprintf("frame %d, playing %v", s.frame+1, s.playing))
	}
}

func demoGlyphMap(c *ui.Context, log func(string)) {
	s := codeState(c)
	if code.GlyphMap(c, "Font", &s.glyph, code.GlyphMapOptions{}).Changed {
		log(fmt.Sprintf("U+%04X", s.glyph))
	}
}

// program is 6502 lines and the CPU after each.
var program = []struct {
	line string
	cpu  code.CPU
}{
	{"RESET", code.CPU{SP: 0xFD, P: 0x24, PC: 0xC000}},
	{"LDA #$3F", code.CPU{A: 0x3F, SP: 0xFD, P: 0x24, PC: 0xC002}},
	{"TAX", code.CPU{A: 0x3F, X: 0x3F, SP: 0xFD, P: 0x24, PC: 0xC003}},
	{"INX", code.CPU{A: 0x3F, X: 0x40, SP: 0xFD, P: 0x24, PC: 0xC004}},
	{"ADC #$C1", code.CPU{A: 0x00, X: 0x40, SP: 0xFD, P: 0x27, PC: 0xC006}},
	{"PHA", code.CPU{A: 0x00, X: 0x40, SP: 0xFC, P: 0x27, PC: 0xC007}},
	{"DEY", code.CPU{A: 0x00, X: 0x40, Y: 0xFF, SP: 0xFC, P: 0xA5, PC: 0xC008}},
}

func demoRegisterView(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := codeState(c)
	p := program[s.step]
	ui.Column(c).Gap(k.D(6)).Children(func() {
		pixel.Text(c, fmt.Sprintf("$%04X  %s", p.cpu.PC, p.line), pixel.TextOptions{Slot: pixel.Accent})
		code.RegisterView(c, "CPU", p.cpu, code.RegisterViewOptions{})
		if input.Button(c, "STEP", input.ButtonOptions{Variant: input.Primary}).Clicked {
			s.step = (s.step + 1) % len(program)
			log(program[s.step].line)
		}
	})
}
