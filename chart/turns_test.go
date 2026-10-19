package chart_test

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chart"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestLegendKeyInAVeilsFirstFrameHidesNothing(t *testing.T) {
	for _, veil := range []bool{false, true} {
		hidden, inert := []bool{false, false}, false
		tt := view(400, 200, newClock(), func(c *ui.Context) {
			s := pixel.Current(c)
			s.Inert = inert
			pixel.Scope(c, s, func() {
				chart.Legend(c, "Key", []chart.Series{{Name: "HP"}, {Name: "MP"}}, hidden, chart.LegendOptions{})
			})
		})
		tt.Key(0, ui.KeyTab)
		inert = veil
		tt.Key(0, ui.KeySpace)
		if hidden[0] == veil {
			t.Errorf("veiled %v, Space left %v", veil, fmt.Sprint(hidden))
		}
	}
}
