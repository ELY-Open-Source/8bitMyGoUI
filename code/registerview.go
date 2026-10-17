package code

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// CPU is a 6502's registers; P holds the flags.
type CPU struct {
	A, X, Y, SP, P uint8
	PC             uint16
}

type RegisterViewOptions struct{}

type regKey struct{}

// seen is the CPU last shown and bits it flipped.
type seen struct {
	cpu, flipped CPU
}

const flags = "NV-BDIZC"

// RegisterView shows registers and flags; the last change lights.
func RegisterView(c *ui.Context, text string, cpu CPU, o RegisterViewOptions) ui.Element {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(1)).Role(ui.RoleGroup).Label(text)
	sn := ui.Local(col, regKey{}, func() seen { return seen{cpu: cpu} })
	if cpu != sn.cpu {
		was := sn.cpu
		sn.cpu, sn.flipped = cpu, CPU{A: was.A ^ cpu.A, X: was.X ^ cpu.X, Y: was.Y ^ cpu.Y, SP: was.SP ^ cpu.SP, P: was.P ^ cpu.P, PC: was.PC ^ cpu.PC}
	}
	fl := sn.flipped
	return col.Children(func() {
		for _, r := range []struct {
			name      string
			v, change uint16
			wide      bool
		}{
			{"A", uint16(cpu.A), uint16(fl.A), false}, {"X", uint16(cpu.X), uint16(fl.X), false}, {"Y", uint16(cpu.Y), uint16(fl.Y), false},
			{"SP", uint16(cpu.SP), uint16(fl.SP), false}, {"PC", cpu.PC, fl.PC, true},
		} {
			digits := 2
			if r.wide {
				digits = 4
			}
			hexed := fmt.Sprintf("$%0*X", digits, r.v)
			ui.Box(c).Size(k.D(120), k.D(pixel.LineHeight)).Role(ui.RoleText).Label(r.name + " " + hexed).Draw(func(p *ui.Painter, rc ui.Rect) {
				g := k.Canvas(p, rc)
				g.Micro(0, 4, r.name, pixel.Faint, 1)
				ink := pixel.Ink
				if r.change != 0 {
					ink = pixel.Accent
				}
				g.Text(12, 0, hexed, ink, 1)
				if !r.wide {
					bits(g, 52, uint8(r.v), uint8(r.change))
				}
			})
		}
		set := []string{}
		for i := range 8 {
			if cpu.P&(0x80>>i) != 0 && flags[i] != '-' {
				set = append(set, string(flags[i]))
			}
		}
		ui.Box(c).Size(k.D(8*9), k.D(14)).Role(ui.RoleText).Label(fmt.Sprintf("P $%02X", cpu.P)).Description(strings.Join(set, " ")).Draw(func(p *ui.Painter, rc ui.Rect) {
			g := k.Canvas(p, rc)
			for i := range 8 {
				x, on, flip := i*9, cpu.P&(0x80>>i) != 0, fl.P&(0x80>>i) != 0
				fill, ink := pixel.Well, pixel.Faint
				if on {
					fill, ink = pixel.Good, pixel.AccentInk
				}
				edge := pixel.Shade
				if flip {
					edge = pixel.Accent
				}
				g.Rect(x, 0, 8, 14, edge)
				g.Rect(x+1, 1, 6, 12, fill)
				g.Micro(x+3-1, 5, string(flags[i]), ink, 1)
			}
		})
	})
}

// bits draws v's eight bits; flipped ones light.
func bits(g pixel.Canvas, x int, v, flipped uint8) {
	for i := range 8 {
		b, ink := "0", pixel.Faint
		if v&(0x80>>i) != 0 {
			b, ink = "1", pixel.Ink
		}
		if flipped&(0x80>>i) != 0 {
			ink = pixel.Accent
		}
		g.Micro(x+i*4+i/4*2, 4, b, ink, 1)
	}
}
