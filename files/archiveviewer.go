package files

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Entry is a file packed in an archive.
type Entry struct {
	Path string
	// Size and Packed are bytes before and after packing.
	Size, Packed int64
}

type ArchiveViewerOptions struct {
	// Height is the rows' in dots; zero means 120.
	Height   int
	Disabled bool
}

type ArchiveViewerResult struct {
	Element *ui.Element
	// Sorted reports a header click this frame.
	Changed, Sorted bool
}

type archiveKey struct{}

// archiveView is the sort and the row chosen in it.
type archiveView struct {
	sort data.Sort
	row  int
}

// ArchiveViewer lists entries and sizes; headers sort them.
func ArchiveViewer(c *ui.Context, text string, entries []Entry, selected *int, o ArchiveViewerOptions) ArchiveViewerResult {
	k := pixel.K(c)
	pick("ArchiveViewer", selected, len(entries))
	var total, packed int64
	for _, e := range entries {
		// Checked against what is left, so sums cannot overflow.
		if e.Path == "" || e.Size < 0 || e.Packed < 0 || e.Size > math.MaxInt64-total || e.Packed > math.MaxInt64-packed {
			panic(fmt.Sprintf("files: ArchiveViewer entry %q of %d, %d bytes", e.Path, e.Size, e.Packed))
		}
		total, packed = total+e.Size, packed+e.Packed
	}
	col := ui.Column(c).Gap(k.D(3))
	st := ui.Local(col, archiveKey{}, func() archiveView { return archiveView{} })
	// The view sorts; the caller's index stays on its entry.
	order := make([]int, len(entries))
	for i := range order {
		order[i] = i
	}
	slices.SortStableFunc(order, func(a, b int) int {
		ea, eb := entries[a], entries[b]
		n := [3]int{cmp.Compare(ea.Path, eb.Path), cmp.Compare(ea.Size, eb.Size), cmp.Compare(ea.Packed, eb.Packed)}[st.sort.Column]
		if st.sort.Desc {
			return -n
		}
		return n
	})
	st.row = slices.Index(order, *selected)
	var res ArchiveViewerResult
	cols := []data.Column{{Title: k.Msg("name"), Width: 112}, {Title: k.Msg("size"), Width: 56}, {Title: k.Msg("packed"), Width: 64}}
	res.Element = col.Children(func() {
		t := data.Table(c, text, cols, len(entries), &st.sort, &st.row, data.TableOptions{Height: o.Height, Disabled: o.Disabled}, func(r, col int) string {
			e := entries[order[r]]
			return [3]string{e.Path, size(e.Size), size(e.Packed)}[col]
		})
		if t.Changed && st.row >= 0 {
			*selected, res.Changed = order[st.row], true
		}
		res.Sorted = t.Sorted
		pixel.Text(c, fmt.Sprintf(k.Msg("archive_total"), len(entries), size(total), size(packed)), pixel.TextOptions{Slot: pixel.Faint})
	})
	if res.Changed {
		slog.Debug("8bit: archive viewer", "label", text, "entry", *selected)
	}
	return res
}
