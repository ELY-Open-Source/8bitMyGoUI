package pixel

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/egoist/mygo/ui"
)

var due atomic.Int64

// Due is the earliest repaint asked for since ResetDue.
func Due() time.Time {
	n := due.Load()
	if n == 0 {
		return time.Time{}
	}
	return time.Unix(0, n)
}

func ResetDue() { due.Store(0) }

func ask(at int64) {
	for {
		cur := due.Load()
		if cur != 0 && cur <= at || due.CompareAndSwap(cur, at) {
			return
		}
	}
}

// Step is the frame at fps, frozen under Reduce Motion.
func (g Canvas) Step(fps int) int {
	checkFPS(fps)
	if g.k.reduce {
		return 0
	}
	return g.Busy(fps)
}

// Busy is the frame at fps, ignoring Reduce Motion.
func (g Canvas) Busy(fps int) int {
	n, d := frameAt(g.now(), fps)
	g.p.After(d)
	ask(g.p.Now().Add(d).UnixNano())
	return n
}

// Tick is the frame at fps, rebuilding the view.
func (k Kit) Tick(c *ui.Context, fps int) int {
	n, d := frameAt(k.now, fps)
	c.After(d)
	ask(c.Now().Add(d).UnixNano())
	return n
}

// Since counts steps from t at fps, up to done.
func (k Kit) Since(c *ui.Context, t time.Time, fps, done int) int {
	if t.IsZero() {
		return done
	}
	n := int(k.now.Sub(t) * time.Duration(fps) / time.Second)
	if n < done {
		k.Tick(c, fps)
	}
	return min(n, done)
}

// frameAt is the frame at t and the next wait.
func frameAt(t time.Time, fps int) (int, time.Duration) {
	checkFPS(fps)
	period := int64(time.Second) / int64(fps)
	ns := t.UnixNano()
	n := ns / period
	return int(n), time.Duration((n+1)*period - ns)
}

func checkFPS(fps int) {
	if fps < 1 || fps > 60 {
		panic(fmt.Sprintf("pixel: fps %d outside 1–60", fps))
	}
}
