package data_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
)

// logs captures debug logs while fn runs.
func logs(fn func()) string {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(old)
	fn()
	return buf.String()
}

func TestStateChangesAreLogged(t *testing.T) {
	sort, row, slot, file := data.Sort{}, 0, 0, "/"
	items := []data.Item{{Sprite: icon.Potion, Name: "Potion"}, {Sprite: icon.Gem, Name: "Gem"}}
	kids := map[string][]string{"/": {"/a"}}
	tt := view(700, 300, newClock(), false, func(c *ui.Context) {
		ui.Row(c).Children(func() {
			data.Table(c, "Roster", []data.Column{{Title: "Name"}}, len(heroes), &sort, &row, data.TableOptions{}, func(r, _ int) string { return heroes[r] })
			data.Inventory(c, "Bag", items, &slot, data.InventoryOptions{})
			data.Outline(c, "Files", []string{"/"}, func(k string) []string { return kids[k] }, func(k string) string { return k }, &file, data.OutlineOptions{})
		})
	})
	out := logs(func() {
		click(t, tt, "Brom")
		tt.Key(0, ui.KeyDown)
		r := find(t, tt, "Bag")
		tt.ClickAt(r.X+4, r.Y+4)
		tt.Key(0, ui.KeyRight)
		click(t, tt, "/")
		tt.Key(0, ui.KeyRight)
	})
	for _, want := range []string{"8bit: table chose", "8bit: inventory chose", "8bit: outline toggled"} {
		if !strings.Contains(out, want) {
			t.Errorf("no %q in the logs:\n%s", want, out)
		}
	}
}
