package tabletop

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TicTacToeOptions struct {
	Disabled bool
}

type TicTacToeResult struct {
	Element *ui.Element
	// Played is the cell you took this frame, else -1.
	Played int
	// Over is the game ending this frame; Restarted, the CONTINUE.
	Over, Restarted bool
	// Winner: 1 is you, 2 the computer, 0 a draw.
	Winner int
}

type tttKey struct{}

// ttt holds the board (1 X, 2 O) and cursor.
type ttt struct {
	board       [9]int
	cur, winner int
	over        bool
}

var rows3 = [8][3]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}

const tcell = 22

// TicTacToe: you are X, the computer O; Enter goes again.
func TicTacToe(c *ui.Context, text string, o TicTacToeOptions) TicTacToeResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(4)).AlignItems(ui.Center).AlignSelf(ui.Start)
	st := ui.Local(col, tttKey{}, func() ttt { return ttt{cur: 4} })
	res := TicTacToeResult{Played: -1}
	res.Element = col.Children(func() {
		shell, base, ps, picked := pixel.Pad(c, text, 3, 9, tcell, tcell, &st.cur, o.Disabled)
		switch {
		case picked < 0:
		case st.over:
			// After the end, a pick is the CONTINUE.
			*st, res.Restarted = ttt{cur: st.cur}, true
			slog.Debug("8bit: tic-tac-toe restart", "label", text)
		case st.board[picked] == 0:
			st.board[picked], res.Played = 1, picked
			reply := -1
			if !st.end() {
				reply = answer(st.board)
				st.board[reply] = 2
				st.end()
			}
			res.Over, res.Winner = st.over, st.winner
			slog.Debug("8bit: tic-tac-toe move", "label", text, "x", picked, "o", reply)
		}
		base.Description([...]string{k.Msg("empty"), "X", "O"}[st.board[st.cur]])
		at, shown := st.cur, st.board
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(0, 0, g.W, g.H, pixel.Lit)
			for i := 1; i < 3; i++ {
				g.Rect(i*tcell-1, 2, 2, 3*tcell-4, pixel.Ink)
				g.Rect(2, i*tcell-1, 3*tcell-4, 2, pixel.Ink)
			}
			for i, m := range shown {
				x, y := i%3*tcell, i/3*tcell
				if i == at && ps.Focus {
					g.Cursor(x+2, y+2, tcell-4, tcell-4)
				}
				switch m {
				case 1:
					for d := range 12 {
						g.Rect(x+5+d, y+5+d, 2, 2, pixel.Bad)
						g.Rect(x+16-d, y+5+d, 2, 2, pixel.Bad)
					}
				case 2:
					g.Circle(x+tcell/2, y+tcell/2, 6, pixel.Info)
					g.Circle(x+tcell/2, y+tcell/2, 5, pixel.Info)
				}
			}
			g.Finish(ps)
		})
		if st.over {
			word := [...]string{k.Msg("draw"), k.Msg("you_win"), k.Msg("game_over")}[st.winner]
			pixel.Text(c, word, pixel.TextOptions{Slot: [...]pixel.Slot{pixel.Ink, pixel.Good, pixel.Bad}[st.winner]})
			pixel.Text(c, k.Msg("continue"), pixel.TextOptions{Slot: pixel.Faint})
		}
	})
	if res.Over {
		slog.Debug("8bit: tic-tac-toe over", "label", text, "winner", res.Winner)
	}
	return res
}

// end checks three in a row, or a full board.
func (t *ttt) end() bool {
	for _, l := range rows3 {
		if m := t.board[l[0]]; m != 0 && m == t.board[l[1]] && m == t.board[l[2]] {
			t.over, t.winner = true, m
			return true
		}
	}
	for _, m := range t.board {
		if m == 0 {
			return false
		}
	}
	t.over, t.winner = true, 0
	return true
}

// answer is O's move: win, block, middle, corners, sides.
func answer(b [9]int) int {
	for _, who := range []int{2, 1} {
		for _, l := range rows3 {
			n, open := 0, -1
			for _, i := range l {
				switch b[i] {
				case who:
					n++
				case 0:
					open = i
				}
			}
			if n == 2 && open >= 0 {
				return open
			}
		}
	}
	for _, i := range []int{4, 0, 2, 6, 8, 1, 3, 5, 7} {
		if b[i] == 0 {
			return i
		}
	}
	panic("tabletop: TicTacToe answered a full board")
}
