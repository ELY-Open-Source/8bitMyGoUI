package files

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type FileDetailsOptions struct {
	// Width is in dots; zero means 220.
	Width int
}

// FileDetails shows a file's icon, name, kind, size and dates.
func FileDetails(c *ui.Context, f File, o FileDetailsOptions) ui.Element {
	k := pixel.K(c)
	check("FileDetails", f)
	return ui.Column(c).Width(k.D(pixel.Or("FileDetails Width", o.Width, 220))).Gap(k.D(6)).Role(ui.RoleGroup).Label(f.Name).Children(func() {
		ui.Row(c).FillWidth().Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
			FileIcon(c, f.Kind, FileIconOptions{Scale: 3})
			pixel.Text(c, f.Name, pixel.TextOptions{Wrap: true, MaxLines: 2}).Grow(1)
		})
		data.DescriptionList(c, [][2]string{
			{k.Msg("kind"), k.Msg(look(f.Kind).msg)},
			{k.Msg("size"), size(f.Size) + " (" + grouped(f.Size) + " B)"},
			{k.Msg("created"), when(f.Created)},
			{k.Msg("modified"), when(f.Modified)},
		}, data.DescriptionListOptions{TermWidth: 64})
	})
}

// when is a date and time, or a dash.
func when(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04")
}

// grouped writes n with commas between thousands.
func grouped(n int64) string {
	s := strconv.FormatInt(n, 10)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}
