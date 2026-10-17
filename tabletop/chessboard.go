package tabletop

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Kind is a chess piece's kind; None, an empty square.
type Kind int

const (
	None Kind = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

var kinds = [...]struct {
	msg string
	sp  *pixel.Sprite
}{{}, {"pawn", icon.Pawn}, {"knight", icon.Knight}, {"bishop", icon.Bishop}, {"rook", icon.Rook}, {"chess_queen", icon.Queen}, {"chess_king", icon.King}}

// Piece is what stands on a square.
type Piece struct {
	Kind  Kind
	White bool
}

type ChessboardOptions struct {
	Disabled bool
}

type ChessboardResult struct {
	Element ui.Element
	// Moved is this frame's move, From one square To another.
	Moved    bool
	From, To int
	// Changed is a pick-up or put-down; Held, the square.
	Changed bool
	Held    int
}

type chessKey struct{}

// chess is the cursor and the held square, else -1.
type chess struct{ cur, held int }

const square = 16

// Chessboard: click a piece, then a square to move to.
func Chessboard(c *ui.Context, text string, board *[64]Piece, o ChessboardOptions) ChessboardResult {
	k := pixel.K(c)
	for i, pc := range board {
		if pc.Kind < None || pc.Kind > King {
			panic(fmt.Sprintf("tabletop: Chessboard square %d has kind %d", i, pc.Kind))
		}
	}
	row := ui.Row(c).Gap(k.D(2)).Role(ui.RoleGroup).Label(text)
	st := ui.Local(row, chessKey{}, func() chess { return chess{held: -1} })
	res := ChessboardResult{From: -1, To: -1, Held: -1}
	res.Element = row.Children(func() {
		ui.Box(c).Size(k.D(5), k.D(8*square)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range 8 {
				g.Micro(1, i*square+6, fmt.Sprint(8-i), pixel.Faint, 1)
			}
		})
		ui.Column(c).Gap(k.D(2)).Children(func() {
			shell, base, ps, picked := pixel.Pad(c, k.Msg("chessboard"), 8, 64, square, square, &st.cur, o.Disabled)
			// A held square the caller emptied is let go.
			if st.held >= 0 && board[st.held].Kind == None {
				st.held = -1
			}
			was := st.held
			switch {
			case picked < 0:
			case picked == st.held:
				st.held = -1
			case st.held < 0 || board[picked].Kind != None && board[picked].White == board[st.held].White:
				// A piece is picked up; its own side swaps it.
				st.held = -1
				if board[picked].Kind != None {
					st.held = picked
				}
			default:
				board[picked], board[st.held] = board[st.held], Piece{}
				res.Moved, res.From, res.To, st.held = true, st.held, picked, -1
			}
			res.Held = st.held
			if st.held != was && !res.Moved {
				res.Changed = true
				slog.Debug("8bit: chess hold", "label", text, "square", st.held)
			}
			base.Description(squareName(st.cur) + " " + pieceName(k, board[st.cur]))
			held, at, shown := st.held, st.cur, *board
			light, dark := shades(k)
			shell.Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for i, pc := range shown {
					x, y := i%8*square, i/8*square
					fill := pixel.Lit
					if (i%8+i/8)%2 == 1 {
						fill = pixel.Shade
					}
					g.Rect(x, y, square, square, fill)
					if i == held {
						g.Rect(x, y, square, square, pixel.AccentLit)
					}
					if i == at && ps.Focus {
						g.Cursor(x, y, square, square)
					}
					if pc.Kind != None {
						piece(g, x+2, y+2, pc, light, dark)
					}
				}
				g.Finish(ps)
			})
			ui.Box(c).Size(k.D(8*square), k.D(5)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for i := range 8 {
					g.Micro(i*square+6, 0, string(rune('A'+i)), pixel.Faint, 1)
				}
			})
		})
	})
	if res.Moved {
		slog.Debug("8bit: chess move", "label", text, "from", squareName(res.From), "to", squareName(res.To))
	}
	return res
}

// piece draws a piece outlined, showing on either square.
func piece(g pixel.Canvas, x, y int, pc Piece, light, dark pixel.Slot) {
	sp, fill, edge := kinds[pc.Kind].sp, dark, light
	if pc.White {
		fill, edge = light, dark
	}
	for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		g.Silhouette(x+d[0], y+d[1], sp, 0, edge, 1)
	}
	g.Sprite(x, y, sp, 0, fill, 1)
}

// squareName is a square as players say it, as E4.
func squareName(i int) string {
	return fmt.Sprintf("%c%d", 'A'+i%8, 8-i/8)
}

// pieceName reads a piece out, as WHITE KNIGHT.
func pieceName(k pixel.Kit, pc Piece) string {
	if pc.Kind == None {
		return k.Msg("empty")
	}
	side := k.Msg("black")
	if pc.White {
		side = k.Msg("white")
	}
	return fmt.Sprintf(k.Msg("piece_name"), side, k.Msg(kinds[pc.Kind].msg))
}
