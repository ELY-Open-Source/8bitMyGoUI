// Package pixel draws 16-bit pixel art for MyGo.
package pixel

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

// Slot names a palette color by its role.
type Slot uint8

const (
	Ink Slot = iota
	Back
	Face
	Lit
	Shade
	Edge
	Faint
	Well
	Accent
	AccentLit
	AccentInk
	Good
	Warn
	Bad
	Info
	Magic
)

const Slots = 16

type Season uint8

const (
	Spring Season = iota
	Summer
	Autumn
	Winter
)

type Palette struct {
	Name, Zh string
	Season   Season
	Night    bool
	Colors   [Slots]ui.Color
}

func (p *Palette) C(s Slot) ui.Color { return p.Colors[s] }

var (
	SpringDay   = palette("Spring Day", "春·昼", Spring, false, "33283A F4EEDC E4DCC6 FFFBEF B3A68C 3A2A3E 6F6358 FFF6E6 B8395C F2A7B9 FFFFFF 4F9A3F D9982A C73A3A 3F7FCB 8C5BC4")
	SpringNight = palette("Spring Night", "春·夜", Spring, true, "F3EEFF 1C1730 2E2850 4C4478 1A1530 E8DFF7 A79DC4 14102A F07EA0 FFC2D2 120E24 7CC66A F2C14E F0605F 6FA8F0 B58AF0")
	SummerDay   = palette("Summer Day", "夏·昼", Summer, false, "15303B F7F3E3 DCEBEE FFFFFF 98B3B8 1D3640 4F6870 F4FCFF 0E6E8C 7FD0E3 FFFFF6 3A9A4A E0A21A D14836 2F7FD0 9A57C8")
	SummerNight = palette("Summer Night", "夏·夜", Summer, true, "EAF3FF 0B1630 172A55 2C4A86 0A1328 DDEBFF 98ACCF 08102A 36C2E0 A8EEFA 050C1E B8E35A F7D04A F26A5A 5FA8FF C58CFF")
	AutumnDay   = palette("Autumn Day", "秋·昼", Autumn, false, "33200F F2E6D0 E6D2B0 FFF6E2 B08F62 3A2416 6E5640 FFF1D8 B94A1C F0A060 FFFFFF 6E8F2E C98A12 B8302A 3F78B0 8E4E9E")
	AutumnNight = palette("Autumn Night", "秋·夜", Autumn, true, "FBEBD8 1E1414 33221E 5A3D30 170F0E F1DCC4 B39A86 140C0B F09A3A FFD08A 0E0808 A8C25A F6C65A F0624E 79A8E0 C88AD8")
	WinterDay   = palette("Winter Day", "冬·昼", Winter, false, "17243A EEF3F8 D8E2EC FFFFFF 98A9BC 1C2A3C 52637A F6FAFF 1F5FA8 9CC8F0 FCFEFF 3C8C6C B98A1E B23A48 2A86A8 7656B8")
	WinterNight = palette("Winter Night", "冬·夜", Winter, true, "EEF7FF 0A1020 18233F 304674 080C1A E2F2FF 9DB2CF 060A16 7FD8F0 CFF4FF 040814 6FE0A8 F5D47A F07A8A 6FB0FF A994F5")
)

var Palettes = []*Palette{SpringDay, SpringNight, SummerDay, SummerNight, AutumnDay, AutumnNight, WinterDay, WinterNight}

func Find(s Season, night bool) *Palette {
	for _, p := range Palettes {
		if p.Season == s && p.Night == night {
			return p
		}
	}
	panic("pixel: no palette for season")
}

// Cycle rotates band by step, for color cycling.
func Cycle(band []Slot, step int) []Slot {
	n := len(band)
	if n == 0 {
		panic("pixel: Cycle needs a band")
	}
	r := (step%n + n) % n
	out := make([]Slot, n)
	for i := range band {
		out[i] = band[(i+r)%n]
	}
	return out
}

func palette(name, zh string, s Season, night bool, hex string) *Palette {
	f := strings.Fields(hex)
	if len(f) != Slots {
		panic("pixel: palette " + name + " needs 16 colors")
	}
	p := &Palette{Name: name, Zh: zh, Season: s, Night: night}
	for i, h := range f {
		p.Colors[i] = ui.Hex("#" + h)
	}
	return p
}
