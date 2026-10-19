package files

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/datetime"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type RecentFilesOptions struct {
	// Width and Height are dots; zeros mean 220, 120.
	Width, Height int
	Disabled      bool
}

type RecentFilesResult struct {
	Element ui.Element
	// Opened is Enter or a double click on a file.
	Changed, Opened bool
}

// RecentFiles lists files with how long since each changed.
func RecentFiles(c *ui.Context, text string, files []File, selected *int, o RecentFilesOptions) RecentFilesResult {
	k := pixel.K(c)
	pick("RecentFiles", selected, len(files))
	for _, f := range files {
		check("RecentFiles", f)
		if f.Modified.IsZero() {
			panic(fmt.Sprintf("files: RecentFiles %q has no time", f.Name))
		}
	}
	var res RecentFilesResult
	res.Element = ui.Column(c).Width(k.D(pixel.Or("RecentFiles Width", o.Width, 220))).Children(func() {
		l := data.List(c, text, len(files), selected, data.ListOptions{Height: o.Height, Disabled: o.Disabled}, func(i int) string { return files[i].Name }, func(i int) {
			row(c, k, files[i], func() { datetime.RelativeTime(c, files[i].Modified, datetime.RelativeTimeOptions{}) })
		})
		l.Element.FillWidth()
		res.Changed, res.Opened = l.Changed, l.Submitted && *selected >= 0
	})
	if res.Opened {
		slog.Debug("8bit: recent files open", "label", text, "file", *selected)
	}
	return res
}
