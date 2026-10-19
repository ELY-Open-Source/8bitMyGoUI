package data_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestCollapsibleKeyInAVeilsFirstFrameOpensNothing(t *testing.T) {
	for _, veil := range []bool{false, true} {
		open, inert := false, false
		tt := underVeil(&inert, func(c *ui.Context) {
			data.Collapsible(c, "Quests", &open, data.CollapsibleOptions{}, func() { pixel.Text(c, "Slay", pixel.TextOptions{}) })
		})
		tt.Key(0, ui.KeyTab)
		inert = veil
		tt.Key(0, ui.KeySpace)
		if open == veil {
			t.Errorf("veiled %v, Space left it open %v", veil, open)
		}
	}
}
