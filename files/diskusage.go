package files

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Use is the room one sort of file takes.
type Use struct {
	Name  string
	Bytes int64
	// Slot colors it; zero picks one by order.
	Slot pixel.Slot
}

type DiskUsageOptions struct {
	// Width is in dots; zero means 200.
	Width int
}

var useSlots = []pixel.Slot{pixel.Info, pixel.Good, pixel.Warn, pixel.Magic, pixel.Bad, pixel.Accent}

// DiskUsage stacks each use in one bar over a legend.
func DiskUsage(c *ui.Context, text string, uses []Use, capacity int64, o DiskUsageOptions) *ui.Element {
	k := pixel.K(c)
	if capacity <= 0 {
		panic(fmt.Sprintf("files: DiskUsage capacity %d", capacity))
	}
	slots, said, used := make([]pixel.Slot, len(uses)), make([]string, len(uses)), int64(0)
	for i, u := range uses {
		// Checked against what is left, so sums cannot overflow.
		if u.Name == "" || u.Bytes < 0 || u.Bytes > capacity-used {
			panic(fmt.Sprintf("files: DiskUsage use %q of %d bytes past %d", u.Name, u.Bytes, capacity))
		}
		used += u.Bytes
		slots[i] = u.Slot
		if slots[i] == 0 {
			slots[i] = useSlots[i%len(useSlots)]
		}
		said[i] = u.Name + " " + size(u.Bytes)
	}
	free := capacity - used
	w := pixel.Or("DiskUsage Width", o.Width, 200)
	return ui.Column(c).Width(k.D(w)).Gap(k.D(4)).Children(func() {
		desc := fmt.Sprintf(k.Msg("used_of"), size(used), size(capacity)) + ", " + strings.Join(append(said, k.Msg("free")+" "+size(free)), ", ")
		ui.Box(c).FillWidth().Height(k.D(12)).Role(ui.RoleMeter).Label(text).Description(desc).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, 12, pixel.Sunken, pixel.Well)
			// Edges come from running sums, so dots never drift.
			iw, sum, from := g.W-4, int64(0), 2
			for i, u := range uses {
				sum += u.Bytes
				to := 2 + int(float64(iw)*float64(sum)/float64(capacity))
				g.Rect(from, 2, to-from, 8, slots[i])
				from = to
			}
		})
		for i, u := range uses {
			legend(c, k, u.Name, u.Bytes, slots[i])
		}
		legend(c, k, k.Msg("free"), free, pixel.Well)
	})
}

// legend is a swatch, a name and a size.
func legend(c *ui.Context, k pixel.Kit, name string, n int64, s pixel.Slot) {
	ui.Row(c).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
		ui.Box(c).Size(k.D(7), k.D(7)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(0, 0, 7, 7, pixel.Edge)
			g.Rect(1, 1, 5, 5, s)
		})
		pixel.Text(c, name, pixel.TextOptions{Wrap: true, MaxLines: 1}).Grow(1)
		pixel.Text(c, size(n), pixel.TextOptions{Slot: pixel.Faint})
	})
}
