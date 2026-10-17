package feedback

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Toast is a short message in the queue Toasts shows.
type Toast struct {
	Message string
	// Action labels an optional button.
	Action string
	// For is how long it shows; zero means three seconds.
	For time.Duration
	// id and shown are set as it first shows.
	id    uint64
	shown time.Time
}

type ToastsOptions struct {
	// Width is each toast's in dots; zero means 200.
	Width    int
	Disabled bool
}

type ToastsResult struct {
	// Element holds the toasts on screen.
	Element ui.Element
	// Acted is the toast whose action was clicked, else nil.
	Acted *Toast
}

type toastsKey struct{}

// Toasts shows three timed messages at the bottom.
func Toasts(c *ui.Context, queue *[]Toast, o ToastsOptions) ToastsResult {
	k := pixel.K(c)
	w := pixel.Or("Toasts Width", o.Width, 200)
	now := k.Now()
	kept := (*queue)[:0]
	for _, t := range *queue {
		if t.For < 0 {
			panic("feedback: Toast For is negative")
		}
		life := t.For
		if life == 0 {
			life = 3 * time.Second
		}
		if t.shown.IsZero() || k.Elapsed(&t.shown) < life {
			kept = append(kept, t)
		}
	}
	*queue = kept
	n := min(3, len(*queue))
	last := ui.Local(c.Root(), toastsKey{}, func() uint64 { return 0 })
	for i := range n {
		if (*queue)[i].id == 0 {
			*last++
			(*queue)[i].id, (*queue)[i].shown = *last, now
		}
	}
	if n > 0 {
		k.Tick(c, 4)
	}
	res := ToastsResult{}
	left := (c.Root().Bounds().W - k.D(w)) / 2
	ui.Overlay(c, func() {
		res.Element = ui.Column(c).Absolute().Left(left).Bottom(k.D(8)).Width(k.D(w)).Gap(k.D(4)).Children(func() {
			for i := range n {
				t := (*queue)[i]
				ui.Row(c).Key(t.id).Padding(k.D(5), k.D(6)).Gap(k.D(6)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(t.Message).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
				}).Children(func() {
					pixel.Text(c, t.Message, pixel.TextOptions{Wrap: true, MaxLines: 2}).Grow(1)
					if t.Action != "" && input.Button(c, t.Action, input.ButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
						res.Acted = &t
					}
				})
			}
		})
	})
	if res.Acted != nil {
		*queue = removeToast(*queue, res.Acted.id)
		slog.Debug("8bit: toast action", "message", res.Acted.Message, "action", res.Acted.Action)
	}
	return res
}

// removeToast drops the toast with id from q.
func removeToast(q []Toast, id uint64) []Toast {
	for i := range q {
		if q[i].id == id {
			return append(q[:i], q[i+1:]...)
		}
	}
	return q
}
