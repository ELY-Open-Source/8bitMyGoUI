package gallery

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// App is the gallery that native and web share.
type App struct {
	router  *ui.Router
	palette int
	locale  pixel.Locale
	started bool
	query   string
	search  pixel.Edit
	logs    map[string][]string
}

// New starts at path; the boot screen shows first.
func New(path string) *App {
	r := ui.NewRouter(path)
	r.Transition = ui.TransitionNone
	return &App{router: r, logs: map[string][]string{}}
}

// Path is the page shown, for the web address bar.
func (a *App) Path() string { return a.router.Location() }

// Go opens path, skipping boot and clearing search.
func (a *App) Go(path string) {
	a.started, a.query = true, ""
	if path != a.router.Location() {
		slog.Debug("8bit: gallery go", "path", path)
		a.router.Push(path)
	}
}

func (a *App) View(c *ui.Context) {
	pixel.Use(c, pixel.Settings{Palette: pixel.Palettes[a.palette], Locale: a.locale})
	if !a.started {
		a.boot(c)
		return
	}
	ui.Column(c).Fill().Children(func() {
		a.topBar(c)
		ui.Row(c).Grow(1).AlignItems(ui.Stretch).Children(func() {
			a.sidebar(c)
			ui.Column(c).Grow(1).Children(func() { a.content(c) })
		})
		a.statusBar(c)
	})
	a.keys(c)
}

// keys leaves single keys to a field being typed in.
func (a *App) keys(c *ui.Context) {
	if a.search.Focused && c.Shortcut(0, ui.KeyEscape) {
		a.query = ""
	}
	if pixel.Typing(c) {
		return
	}
	switch {
	case c.Shortcut(0, ui.KeySlash):
		a.search.Input.Focus()
	case c.Shortcut(0, ui.KeyP):
		a.setPalette(a.palette + 1)
	case c.Shortcut(ui.Shift, ui.KeyP):
		a.setPalette(a.palette - 1)
	case c.Shortcut(0, ui.KeyBracketLeft):
		a.step(-1)
	case c.Shortcut(0, ui.KeyBracketRight):
		a.step(1)
	}
}

func (a *App) setPalette(i int) {
	n := len(pixel.Palettes)
	a.palette = (i%n + n) % n
	slog.Debug("8bit: gallery palette", "name", pixel.Palettes[a.palette].Name)
}

// step moves to the previous or next entry.
func (a *App) step(d int) {
	all := Entries()
	for i, e := range all {
		if e.Path() == a.router.Path() {
			a.Go(all[((i+d)%len(all)+len(all))%len(all)].Path())
			return
		}
	}
}

// logger records a demo's events for the log panel.
func (a *App) logger(e Entry) func(string) {
	return func(ev string) {
		slog.Debug("8bit: demo", "entry", e.Slug(), "event", ev)
		l := append(a.logs[e.Slug()], ev)
		a.logs[e.Slug()] = l[max(0, len(l)-8):]
	}
}

func find(slug string) (Entry, bool) {
	for _, e := range Entries() {
		if e.Slug() == slug {
			return e, true
		}
	}
	return Entry{}, false
}

func matches(q string) []Entry {
	q = strings.ToLower(strings.TrimSpace(q))
	var out []Entry
	for _, e := range Entries() {
		if strings.Contains(strings.ToLower(e.Name+" "+e.Zh+" "+e.Slug()), q) {
			out = append(out, e)
		}
	}
	return out
}
