//go:build js && wasm

// Command web runs the gallery in a browser canvas.
package main

import (
	"fmt"
	"log/slog"
	"runtime/debug"
	"syscall/js"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/gallery"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// dot is DIPs per dot; each dot is one pixel.
const dot = gallery.Dot

type host struct {
	app    *gallery.App
	tt     *ui.Tester
	w, h   int
	frames int
	opened int
	over   bool
	dead   bool
	began  time.Time
}

func main() {
	h := &host{app: gallery.New("/")}
	api := js.Global().Get("Object").New()
	export := func(name string, fn func(args []js.Value) any) {
		api.Set(name, js.FuncOf(func(_ js.Value, args []js.Value) any { return h.guard(name, args, fn) }))
	}
	export("start", h.start)
	export("resize", func(a []js.Value) any { h.size(a[0].Int(), a[1].Int()); return h.out() })
	export("pointer", h.pointer)
	export("wheel", func(a []js.Value) any {
		h.tt.Scroll(dip(a[0]), dip(a[1]), dip(a[2]), dip(a[3]))
		return h.out()
	})
	export("key", h.key)
	export("text", func(a []js.Value) any { h.tt.Type(a[0].String()); return h.out() })
	export("compose", func(a []js.Value) any { h.tt.Compose(a[0].String(), a[1].Int()); return h.out() })
	export("paste", func(a []js.Value) any {
		h.tt.SetClipboard(a[0].String())
		h.tt.Command("paste")
		h.tt.SetClipboard("")
		return h.out()
	})
	export("focus", func(a []js.Value) any { h.tt.SetFocused(a[0].Bool()); return h.out() })
	export("go", func(a []js.Value) any { h.app.Go(a[0].String()); h.tt.Frame(); return h.out() })
	export("frame", func([]js.Value) any { h.tt.Frame(); return h.out() })
	export("stats", h.stats)
	export("texts", func([]js.Value) any { return toJS(h.tt.Texts()) })
	export("find", func(a []js.Value) any {
		r, ok := h.tt.Find(a[0].String())
		if !ok {
			return nil
		}
		return []any{float64(r.X / dot), float64(r.Y / dot), float64(r.W / dot), float64(r.H / dot)}
	})
	js.Global().Set("__8bitGo", api)
	slog.Info("8bit: wasm ready")
	select {}
}

// guard runs a JS call; a panic shows GAME OVER.
func (h *host) guard(name string, args []js.Value, fn func([]js.Value) any) (res any) {
	if h.dead {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			slog.Error("8bit: crash", "call", name, "panic", r, "stack", string(debug.Stack()))
			res = h.crash(fmt.Sprint(r))
		}
	}()
	// An empty clipboard after a call means nothing was copied.
	if h.tt != nil {
		h.tt.SetClipboard("")
	}
	pixel.ResetDue()
	h.began = time.Now()
	return fn(args)
}

// crash shows GAME OVER; a second panic ends the page.
func (h *host) crash(why string) (res any) {
	if h.over || h.w == 0 {
		h.dead = true
		js.Global().Call("__8bitCrash", why)
		return nil
	}
	h.over = true
	defer func() {
		if r := recover(); r != nil {
			res = h.crash(fmt.Sprintf("%s\n\nGAME OVER failed: %v", why, r))
		}
	}()
	h.tt = ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.WinterNight, Dot: dot})
		ui.Column(c).Fill().Center().Children(func() {
			if feedback.ErrorState(c, why, feedback.ErrorStateOptions{}).Continue {
				js.Global().Get("location").Call("reload")
			}
		})
	}, h.w*dot, h.h*dot)
	h.tt.SetScale(1.0 / dot)
	return h.out()
}

// start sizes the window and opens any deep link.
func (h *host) start(a []js.Value) any {
	h.w, h.h = a[0].Int(), a[1].Int()
	if path := a[2].String(); path != "" {
		h.app.Go(path)
	}
	h.tt = ui.NewTester(h.app.View, h.w*dot, h.h*dot)
	h.tt.SetScale(1.0 / dot)
	slog.Info("8bit: start", "dots", fmt.Sprintf("%dx%d", h.w, h.h), "path", h.app.Path())
	return h.out()
}

func (h *host) size(w, ht int) {
	h.w, h.h = w, ht
	h.tt.SetSize(w*dot, ht*dot)
	slog.Info("8bit: resize", "dots", fmt.Sprintf("%dx%d", w, ht))
}

func (h *host) pointer(a []js.Value) any {
	x, y := dip(a[1]), dip(a[2])
	switch kind := a[0].String(); kind {
	case "move":
		h.tt.Move(x, y)
	case "down":
		h.tt.Press(x, y)
	case "up":
		h.tt.Release(x, y)
	case "right":
		h.tt.RightClickAt(x, y)
	default:
		panic("web: unknown pointer kind " + kind)
	}
	return h.out()
}

func (h *host) key(a []js.Value) any {
	k, ok := keys[a[0].String()]
	if !ok {
		return nil
	}
	h.tt.Key(ui.Modifiers(a[1].Int()), k)
	return h.out()
}

// out returns the frame and requests for the page.
func (h *host) out() any {
	img := h.tt.Image()
	h.frames++
	pix := js.Global().Get("Uint8ClampedArray").New(len(img.Pix))
	js.CopyBytesToJS(pix, img.Pix)
	// took is Go's share; the page adds its blit.
	took := float64(time.Since(h.began).Microseconds()) / 1000
	res := map[string]any{"w": img.Rect.Dx(), "h": img.Rect.Dy(), "pix": pix, "path": h.app.Path(), "cursor": int(h.tt.Cursor()), "took": took}
	if due := pixel.Due(); !due.IsZero() {
		res["due"] = float64(time.Until(due).Milliseconds())
	}
	if r, ok := h.tt.TextCaret(); ok {
		res["caret"] = []any{float64(r.X / dot), float64(r.Y / dot), float64(r.H / dot)}
	}
	if c := h.tt.Clipboard(); c != "" {
		res["copied"] = c
	}
	if urls := h.tt.OpenedURLs(); len(urls) > h.opened {
		res["open"] = toJS(urls[h.opened:])
		h.opened = len(urls)
	}
	return res
}

func (h *host) stats([]js.Value) any {
	return map[string]any{"frames": h.frames, "w": h.w, "h": h.h, "path": h.app.Path()}
}

func dip(v js.Value) float32 { return float32(v.Float() * dot) }

func toJS(s []string) any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}
