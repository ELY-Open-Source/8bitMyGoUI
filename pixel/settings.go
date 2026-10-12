package pixel

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"
)

type Locale uint8

const (
	En Locale = iota
	Zh
)

type Settings struct {
	Palette *Palette
	// Dot is DIPs per dot, 1–4; zero means 2.
	Dot    int
	Locale Locale
	// Now replaces the frame clock, for tests.
	Now func() time.Time
	// Inert controls take no input, as under a busy veil.
	Inert bool
}

// Kit is one frame's settings, captured by components.
type Kit struct {
	pal    *Palette
	dot    float32
	locale Locale
	now    time.Time
	clock  func() time.Time
	reduce bool
}

func (k Kit) Palette() *Palette { return k.pal }
func (k Kit) C(s Slot) ui.Color { return k.pal.Colors[s] }
func (k Kit) Dot() float32      { return k.dot }
func (k Kit) D(n int) float32   { return float32(n) * k.dot }
func (k Kit) Locale() Locale    { return k.locale }
func (k Kit) Now() time.Time    { return k.now }
func (k Kit) Reduce() bool      { return k.reduce }

type frame struct {
	cur    Settings
	at     time.Time
	typing bool
}

type frameKey struct{}

// Use must run once per frame, before any component.
func Use(c *ui.Context, s Settings) {
	s = s.valid()
	f := frameOf(c)
	if f.cur.Palette != s.Palette {
		slog.Debug("8bit: palette", "name", s.Palette.Name)
	}
	f.cur, f.at, f.typing = s, c.Now(), false
	c.SetTheme(theme(s.Palette))
}

// Typing reports a focused Editor built so far this frame.
func Typing(c *ui.Context) bool { return current(c).typing }

// Current returns the settings in force, for Scope to change.
func Current(c *ui.Context) Settings { return current(c).cur }

// Scope builds fn under other settings.
func Scope(c *ui.Context, s Settings, fn func()) {
	f := current(c)
	prev := f.cur
	f.cur = s.valid()
	fn()
	f.cur = prev
}

// Lazy keeps this scope for rows MyGo builds during layout.
func Lazy[T any](c *ui.Context, fn func(T)) func(T) {
	s := current(c).cur
	return func(v T) { Scope(c, s, func() { fn(v) }) }
}

func K(c *ui.Context) Kit {
	s := current(c).cur
	now := c.Now()
	if s.Now != nil {
		now = s.Now()
	}
	return Kit{pal: s.Palette, dot: float32(s.Dot), locale: s.Locale, now: now, clock: s.Now, reduce: c.Preferences().ReduceMotion}
}

func frameOf(c *ui.Context) *frame {
	return ui.Local(c.Root(), frameKey{}, func() frame { return frame{} })
}

func current(c *ui.Context) *frame {
	f := frameOf(c)
	if !f.at.Equal(c.Now()) {
		panic("pixel: Use was not called in this frame")
	}
	return f
}

func (s Settings) valid() Settings {
	if s.Palette == nil {
		panic("pixel: Settings.Palette is nil")
	}
	if s.Dot == 0 {
		s.Dot = 2
	}
	if s.Dot < 1 || s.Dot > 4 {
		panic(fmt.Sprintf("pixel: Settings.Dot %d outside 1–4", s.Dot))
	}
	if s.Locale > Zh {
		panic(fmt.Sprintf("pixel: unknown Locale %d", s.Locale))
	}
	return s
}

// theme keeps MyGo's drawing in the palette; row highlights clear.
func theme(p *Palette) *ui.Theme {
	t := ui.LightTheme()
	c := p.Colors
	t.Dark = p.Night
	t.Background, t.Surface, t.SurfaceHover, t.SurfacePressed = c[Back], c[Face], ui.Transparent, c[Shade]
	t.Border, t.Text, t.TextMuted = c[Edge], c[Ink], c[Faint]
	t.Accent, t.AccentHover, t.AccentPressed, t.AccentText = ui.Transparent, c[AccentLit], c[Accent], c[AccentInk]
	t.Danger, t.Warning, t.Success = c[Bad], c[Warn], c[Good]
	t.Selection, t.Focus, t.Scrollbar, t.Radius = c[AccentLit], c[Accent], ui.Transparent, 0
	return t
}

// Or is v, or def for zero; negatives panic.
func Or(name string, v, def int) int {
	if v < 0 {
		panic(fmt.Sprintf("pixel: %s %d is negative", name, v))
	}
	if v == 0 {
		return def
	}
	return v
}
