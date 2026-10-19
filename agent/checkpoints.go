package agent

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/datetime"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Checkpoint is a save point to go back to.
type Checkpoint struct {
	Name string
	At   time.Time
}

type CheckpointsOptions struct {
	// Width is in dots; zero means 240.
	Width    int
	Disabled bool
}

type CheckpointsResult struct {
	Element ui.Element
	// Restored is the checkpoint chosen this frame, else -1.
	Restored int
}

// Checkpoints lists save points, the newest aglow; RESTORE goes back.
func Checkpoints(c *ui.Context, text string, points []Checkpoint, o CheckpointsOptions) CheckpointsResult {
	k := pixel.K(c)
	newest, seen := -1, map[string]bool{}
	for i, p := range points {
		if p.Name == "" || p.At.IsZero() || seen[p.Name] {
			panic(fmt.Sprintf("agent: Checkpoints point %d %q is blank, undated or twice", i, p.Name))
		}
		seen[p.Name] = true
		if newest < 0 || points[newest].At.Before(p.At) {
			newest = i
		}
	}
	w := pixel.Or("Checkpoints Width", o.Width, 240)
	// Names and times need room beside the crystal and RESTORE.
	if need := 12 + 4 + 96 + 4 + pixel.Measure(k.Msg("restore")) + 8; w < need {
		panic(fmt.Sprintf("agent: Checkpoints Width %d is under %d", w, need))
	}
	res := CheckpointsResult{Restored: -1}
	col := ui.Column(c).Width(k.D(w)).Gap(k.D(4)).Role(ui.RoleGroup).Label(text)
	if o.Disabled {
		col.Description(k.Msg("disabled"))
	}
	res.Element = col.Children(func() {
		if len(points) == 0 {
			pixel.Text(c, k.Msg("empty"), pixel.TextOptions{Slot: pixel.Faint})
		}
		for i, p := range points {
			// Keyed by name: points come and go.
			ui.Box(c).Key(p.Name).FillWidth().Children(func() {
				ui.Row(c).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
					glow := i == newest
					ui.Box(c).Size(k.D(12), k.D(12)).Role(ui.RoleNone).Draw(func(pt *ui.Painter, r ui.Rect) {
						g := k.Canvas(pt, r)
						if !glow {
							g.Sprite(0, 0, icon.Crystal, 0, pixel.Info, 1)
							return
						}
						// The newest save point shimmers in Magic and Accent.
						g.Sprite(0, 0, icon.Crystal, 0, pixel.Cycle([]pixel.Slot{pixel.Magic, pixel.Accent}, g.Step(4))[0], 1)
					})
					ui.Column(c).Grow(1).Children(func() {
						pixel.Text(c, p.Name, pixel.TextOptions{Wrap: true, MaxLines: 1})
						datetime.RelativeTime(c, p.At, datetime.RelativeTimeOptions{})
					})
					if input.Button(c, k.Msg("restore"), input.ButtonOptions{Size: input.Small, Label: k.Msg("restore") + " " + p.Name, Disabled: o.Disabled}).Clicked {
						res.Restored = i
					}
				})
			})
		}
	})
	if res.Restored >= 0 {
		slog.Debug("8bit: checkpoint restored", "label", text, "name", points[res.Restored].Name)
	}
	return res
}
