package desk

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TaskbarOptions struct {
	// Programs fill the START menu; none leaves START out.
	Programs []Program
	// Tray icons sit by the clock.
	Tray []TrayIcon
	// Width is in dots; zero means 320.
	Width    int
	Disabled bool
}

type TaskbarResult struct {
	Element *ui.Element
	// Changed is a window button pressed; Launched a program.
	Changed  bool
	Launched string
	// Tray is the tray icon pressed this frame, else -1.
	Tray int
}

type barKey struct{}

type slideKey struct{}

const taskW = 72

// Taskbar: START, a button per window, tray and clock.
func Taskbar(c *ui.Context, text string, windows []string, active *int, o TaskbarOptions) TaskbarResult {
	k := pixel.K(c)
	if *active < -1 || *active >= len(windows) {
		panic(fmt.Sprintf("desk: Taskbar active %d of %d windows", *active, len(windows)))
	}
	w := pixel.Or("Taskbar Width", o.Width, 320)
	row := ui.Row(c).Width(k.D(w)).Height(k.D(pixel.LineHeight + 10)).Padding(k.D(3)).Gap(k.D(3)).AlignItems(ui.Center).Role(ui.RoleToolbar).Label(text)
	cur := ui.Local(row, barKey{}, func() int { return max(0, *active) })
	first := ui.Local(row, slideKey{}, func() int { return 0 })
	res := TaskbarResult{Tray: -1}
	res.Element = row.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Face)
		g.HLine(0, 0, g.W, pixel.Lit)
	}).Children(func() {
		if len(o.Programs) > 0 {
			res.Launched = StartMenu(c, k.Msg("start"), o.Programs, StartMenuOptions{Disabled: o.Disabled}).Launched
		}
		// Window buttons share the room; arrows reach any past it.
		room := ui.Box(c).Key("windows").Grow(1).Height(k.D(pixel.LineHeight + 4)).Clip()
		free, check := pixel.WrapWidth(c, room)
		room.Draw(func(p *ui.Painter, r ui.Rect) { check(r) }).Children(func() {
			if len(windows) == 0 {
				return
			}
			bw := max(24, min(taskW, free/len(windows)))
			*cur = min(*cur, len(windows)-1)
			was := *active
			shell, base, st, picked := pixel.Pad(c, k.Msg("open_windows"), len(windows), len(windows), bw, pixel.LineHeight+4, cur, o.Disabled)
			// Pressing the active window's button minimizes it.
			if picked >= 0 {
				*active, res.Changed = picked, true
				if picked == was {
					*active = -1
				}
			}
			// Past the room, the buttons slide so the cursor shows.
			fit := max(1, free/bw)
			*first = max(0, min(*first, *cur, len(windows)-fit))
			if *cur >= *first+fit {
				*first = *cur - fit + 1
			}
			shell.Absolute().Left(-k.D(*first * bw)).Top(0)
			at, on := *cur, *active
			base.Description(windows[at])
			shell.Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for i, name := range windows {
					x, style := i*bw, pixel.Raised
					if i == on {
						style = pixel.Sunken
					}
					g.Frame(x+1, 0, bw-2, g.H, style, pixel.Face)
					if i == on {
						g.Checker(x+3, 2, bw-6, g.H-4, pixel.Lit)
					}
					if i == at && st.Focus {
						g.Cursor(x+3, 3, bw-6, g.H-6)
					}
					g.Text(x+5, 2, pixel.Wrap(name, bw-10, 1)[0], pixel.Ink, 1)
				}
				g.Finish(st)
			})
		})
		// Keyed: the tray keeps focus as START comes and goes.
		ui.Row(c).Key("tray").Gap(k.D(4)).Padding(k.D(2), k.D(4)).AlignItems(ui.Center).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Face)
		}).Children(func() {
			if len(o.Tray) > 0 {
				res.Tray = TrayIcons(c, k.Msg("tray"), o.Tray, TrayIconsOptions{Disabled: o.Disabled}).Clicked
			}
			now := k.Now()
			k.At(c, now.Truncate(time.Minute).Add(time.Minute))
			pixel.Text(c, now.Format("15:04"), pixel.TextOptions{})
		})
	})
	if res.Changed {
		slog.Debug("8bit: taskbar window", "label", text, "active", *active)
	}
	return res
}
