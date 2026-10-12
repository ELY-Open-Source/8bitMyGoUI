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

// At builds the view again at t.
func (k Kit) At(c *ui.Context, t time.Time) {
	d := max(0, t.Sub(k.now))
	c.After(d)
	ask(c.Now().Add(d).UnixNano())
}

// Elapsed: time since *t; clock steps back, *t becomes now.
func (k Kit) Elapsed(t *time.Time) time.Duration {
	if k.now.Before(*t) {
		*t = k.now
	}
	return k.now.Sub(*t)
}

// Since counts steps from t at fps, up to done.
func (k Kit) Since(c *ui.Context, t time.Time, fps, done int) int {
	checkFPS(fps)
	if t.IsZero() {
		return done
	}
	// A clock stepping back holds at zero; division cannot overflow.
	n := 0
	if d := k.now.Sub(t); d > 0 {
		n = int(min(d/(time.Second/time.Duration(fps)), time.Duration(done)))
	}
	if n < done {
		k.Tick(c, fps)
	}
	return min(n, done)
}

type towardKey struct{}

// heading is a Toward's target, last turn and step then.
type heading struct {
	on   bool
	at   time.Time
	from int
}

// Toward steps 0–n toward on, turning from where it is.
func (k Kit) Toward(c *ui.Context, e *ui.Element, on bool, fps, n int) int {
	checkFPS(fps)
	if n < 1 {
		panic(fmt.Sprintf("pixel: Toward needs steps, got %d", n))
	}
	end := func(on bool) int {
		if on {
			return n
		}
		return 0
	}
	// First shown, it already sits where on puts it.
	s := ui.Local(e, towardKey{}, func() heading { return heading{on: on, from: end(on)} })
	// Reduce Motion settles at the end, so nothing resumes later.
	if k.reduce {
		*s = heading{on: on, from: end(on)}
		return s.from
	}
	shown := func() int {
		if s.on {
			return s.from + k.Since(c, s.at, fps, n-s.from)
		}
		return s.from - k.Since(c, s.at, fps, s.from)
	}
	if on != s.on {
		s.from, s.on, s.at = shown(), on, k.now
	}
	return shown()
}

type shotKey struct{}

// shot is a one-shot's cue, start, and spent mark.
type shot struct {
	cue   int
	at    time.Time
	spent bool
}

// Shot steps a cued one-shot; Reduce Motion spends it whole.
func (k Kit) Shot(c *ui.Context, e *ui.Element, cue, fps, done int) int {
	checkFPS(fps)
	// Before any change the zero start reads as done.
	s := ui.Local(e, shotKey{}, func() shot { return shot{cue: cue} })
	if cue != s.cue {
		s.cue, s.at, s.spent = cue, k.now, false
	}
	if k.reduce {
		s.spent = true
	}
	if s.spent {
		return done
	}
	return k.Since(c, s.at, fps, done)
}

type cueKey struct{}

// cued is a cue's last value and when it changed.
type cued struct {
	n  int
	at time.Time
}

// Cue is when n last changed on e, zero before.
func (k Kit) Cue(e *ui.Element, n int) time.Time {
	s := ui.Local(e, cueKey{}, func() cued { return cued{n: n} })
	if s.n != n {
		s.n, s.at = n, k.now
	}
	return s.at
}

// frameAt is the frame at t and the next wait.
func frameAt(t time.Time, fps int) (int, time.Duration) {
	checkFPS(fps)
	period := int64(time.Second) / int64(fps)
	ns := t.UnixNano()
	n := ns / period
	return int(n), time.Duration((n+1)*period - ns)
}

// checkFPS holds every clock to S5's 4–15 fps.
func checkFPS(fps int) {
	if fps < 4 || fps > 15 {
		panic(fmt.Sprintf("pixel: fps %d outside 4–15", fps))
	}
}
