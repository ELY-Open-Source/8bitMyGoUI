package gallery

import (
	"strconv"
	"strings"
	"testing"
)

func TestTabletopDemosLogTheirActions(t *testing.T) {
	// Focusing, or a move that ends nothing, logs nothing.
	quiet := map[string]bool{"key:Right": true}
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/486-dice", []string{"ROLL"}}, {"/c/487-dicetray", []string{"ROLL ALL"}},
		{"/c/489-cardhand", []string{"Hand", "DEAL"}}, {"/c/490-deck", []string{"Draw pile"}},
		{"/c/491-chessboard", []string{"square:52", "NEW GAME"}}, {"/c/492-goboard", []string{"Goban"}},
		{"/c/495-roulette", []string{"SPIN"}}, {"/c/496-coinflip", []string{"FLIP"}},
		{"/c/497-bingocard", []string{"Card", "key:Right", "key:Enter"}}, {"/c/498-tictactoe", []string{"Grid"}},
		{"/c/499-scoretrack", []string{"TAKE A TURN"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			if k, ok := strings.CutPrefix(step, "key:"); ok {
				tt.Key(0, letters[k])
			} else if sq, ok := strings.CutPrefix(step, "square:"); ok {
				// A chess square: 16 dots, two pixels a dot.
				r, ok := tt.Find("CHESSBOARD")
				i, err := strconv.Atoi(sq)
				if !ok || err != nil {
					t.Fatalf("no board, or square %q: %v", sq, err)
				}
				tt.ClickAt(r.X+float32(32*(i%8)+16), r.Y+float32(32*(i/8)+16))
			} else {
				click(t, tt, step)
			}
			if len(a.logs[slug]) <= before && !quiet[step] && step != "Card" {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}
