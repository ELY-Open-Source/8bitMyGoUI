package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Binding is an action and the key that does it.
type Binding struct {
	Action string
	Key    ui.Key
}

// caps hold every key but Escape, which cancels.
var caps = map[ui.Key]string{
	ui.KeyA: "A", ui.KeyB: "B", ui.KeyC: "C", ui.KeyD: "D", ui.KeyE: "E", ui.KeyF: "F", ui.KeyG: "G", ui.KeyH: "H", ui.KeyI: "I",
	ui.KeyJ: "J", ui.KeyK: "K", ui.KeyL: "L", ui.KeyM: "M", ui.KeyN: "N", ui.KeyO: "O", ui.KeyP: "P", ui.KeyQ: "Q", ui.KeyR: "R",
	ui.KeyS: "S", ui.KeyT: "T", ui.KeyU: "U", ui.KeyV: "V", ui.KeyW: "W", ui.KeyX: "X", ui.KeyY: "Y", ui.KeyZ: "Z",
	ui.Key0: "0", ui.Key1: "1", ui.Key2: "2", ui.Key3: "3", ui.Key4: "4", ui.Key5: "5", ui.Key6: "6", ui.Key7: "7", ui.Key8: "8", ui.Key9: "9",
	ui.KeyUp: "↑", ui.KeyDown: "↓", ui.KeyLeft: "←", ui.KeyRight: "→",
	ui.KeySpace: "SPACE", ui.KeyEnter: "ENTER", ui.KeyTab: "TAB", ui.KeyBackspace: "BKSP",
	ui.KeyDelete: "DEL", ui.KeyInsert: "INS", ui.KeyHome: "HOME", ui.KeyEnd: "END", ui.KeyPageUp: "PGUP", ui.KeyPageDown: "PGDN",
	ui.KeyF1: "F1", ui.KeyF2: "F2", ui.KeyF3: "F3", ui.KeyF4: "F4", ui.KeyF5: "F5", ui.KeyF6: "F6",
	ui.KeyF7: "F7", ui.KeyF8: "F8", ui.KeyF9: "F9", ui.KeyF10: "F10", ui.KeyF11: "F11", ui.KeyF12: "F12",
	ui.KeyMinus: "-", ui.KeyEqual: "=", ui.KeyComma: ",", ui.KeyPeriod: ".", ui.KeySlash: "/", ui.KeySemicolon: ";",
	ui.KeyQuote: "'", ui.KeyBracketLeft: "[", ui.KeyBracketRight: "]", ui.KeyBackslash: "\\", ui.KeyBackquote: "`",
	ui.KeyContextMenu: "MENU", ui.KeyBack: "BACK", ui.KeyForward: "FWD",
}

type KeyBindingsOptions struct {
	// Width is in dots; zero means 200.
	Width    int
	Disabled bool
}

type KeyBindingsResult struct {
	Element *ui.Element
	Changed bool
}

type bindKey struct{}

// capture is the row in hand and whether it waits.
type capture struct {
	cur     int
	waiting bool
}

// KeyBindings lists actions with keys; Enter, then a key, rebinds.
func KeyBindings(c *ui.Context, text string, binds []Binding, o KeyBindingsOptions) KeyBindingsResult {
	k := pixel.K(c)
	used := map[ui.Key]int{}
	for _, b := range binds {
		if _, ok := caps[b.Key]; !ok {
			panic(fmt.Sprintf("account: KeyBindings %q has a key without a cap", b.Action))
		}
		used[b.Key]++
	}
	if len(binds) == 0 {
		panic("account: KeyBindings has no bindings")
	}
	w, rh := pixel.Or("KeyBindings Width", o.Width, 200), pixel.LineHeight+4
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleList).Label(text).Disabled(o.Disabled)
	})
	kb := ui.Local(base, bindKey{}, func() capture { return capture{} })
	kb.cur = min(kb.cur, len(binds)-1)
	res := KeyBindingsResult{Element: shell}
	switch {
	case base.IsDisabled():
		kb.waiting = false
	case kb.waiting:
		// Waiting, every key binds, arrows too; Escape gives up.
		if base.Shortcut(0, ui.KeyEscape) {
			kb.waiting = false
		}
		for key := range caps {
			if base.Shortcut(0, key) && kb.waiting {
				binds[kb.cur].Key, kb.waiting, res.Changed = key, false, true
				slog.Debug("8bit: key bindings", "label", text, "action", binds[kb.cur].Action, "key", caps[key])
			}
		}
	default:
		if base.Shortcut(0, ui.KeyUp) {
			kb.cur = max(0, kb.cur-1)
		}
		if base.Shortcut(0, ui.KeyDown) {
			kb.cur = min(len(binds)-1, kb.cur+1)
		}
		if base.Shortcut(0, ui.KeyEnter) || base.Shortcut(0, ui.KeySpace) {
			kb.waiting = true
		}
		if base.Clicked() {
			if _, y, _ := base.PointerPosition(); int(y/k.Dot())/rh < len(binds) {
				kb.cur, kb.waiting = int(y/k.Dot())/rh, true
			}
		}
	}
	b := binds[kb.cur]
	say := b.Action + ", " + caps[b.Key]
	switch {
	case kb.waiting:
		say = b.Action + ", " + k.Msg("press_key")
	case used[b.Key] > 1:
		say += ", " + k.Msg("conflict")
	}
	base.Description(say)
	cur, waiting := kb.cur, kb.waiting
	shell.Size(k.D(w), k.D(len(binds)*rh)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, b := range binds {
			y := i * rh
			if i == cur && st.Focus {
				g.Cursor(0, y, w, rh)
			}
			g.Text(4, y+2, b.Action, pixel.Ink, 1)
			cap, ink, edge := caps[b.Key], pixel.Ink, pixel.Edge
			switch {
			case i == cur && waiting:
				cap, ink = k.Msg("press_key"), pixel.Accent
			case used[b.Key] > 1:
				ink, edge = pixel.Bad, pixel.Bad
			}
			cw := pixel.Measure(cap) + 8
			g.Frame(w-cw-2, y+1, cw, rh-2, pixel.Raised, pixel.Face)
			if edge == pixel.Bad {
				g.Rect(w-cw-2, y+rh-2, cw, 1, pixel.Bad)
			}
			g.Text(w-cw+2, y+2, cap, ink, 1)
		}
		g.Finish(st)
	})
	return res
}
