package diagram

import (
	"errors"
	"fmt"
	"log/slog"

	pele "github.com/Broderick-Westrope/pele-go"
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/feedback"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type MermaidOptions struct {
	// Label names the diagram; empty takes the source's title.
	Label string
}

type mermaidKey struct{}

// drawn is a scene, or its error, and their inputs.
type drawn struct {
	source string
	colors [pixel.Slots]ui.Color
	sc     *scene
	err    error
	// fit holds the scene's cells by shrink factor.
	fit map[int][]cell
}

// Mermaid draws Mermaid source on dots, shrunk to fit.
func Mermaid(c *ui.Context, source string, o MermaidOptions) ui.Element {
	k := pixel.K(c)
	pal := k.Palette()
	root := ui.Column(c).FillWidth()
	d := ui.Local(root, mermaidKey{}, func() drawn { return drawn{} })
	if d.fit == nil || d.source != source || d.colors != pal.Colors {
		*d = drawn{source: source, colors: pal.Colors, fit: map[int][]cell{}}
		d.sc, d.err = render(source, tokens(pal), &pal.Colors)
		d.report()
	}
	return root.Children(func() {
		if d.err != nil {
			body := d.err.Error()
			var pe *pele.Error
			if errors.As(d.err, &pe) && pe.Line > 0 {
				body = fmt.Sprintf("%d:%d %s", pe.Line, pe.Column, body)
			}
			feedback.Alert(c, k.Msg("diagram_failed"), body, feedback.AlertOptions{Kind: feedback.Danger})
			return
		}
		label := o.Label
		if label == "" {
			label = d.sc.title
		}
		if label == "" {
			label = k.Msg("diagram")
		}
		ui.Box(c).FillWidth().MaxWidth(k.D(d.sc.w)).AspectRatio(float32(d.sc.w) / float32(d.sc.h)).
			Role(ui.RoleImage).Label(label).Draw(func(p *ui.Painter, r ui.Rect) {
			// A hair of float error must not cost a dot.
			w, h := int(r.W/k.Dot()+0.01), int(r.H/k.Dot()+0.01)
			if w < 1 || h < 1 {
				return
			}
			g := k.Canvas(p, r)
			for _, cl := range d.cells(max(ceil(d.sc.w, w), ceil(d.sc.h, h))) {
				g.Rect(cl.x, cl.y, cl.w, cl.h, cl.slot)
			}
		})
	})
}

func ceil(a, b int) int { return (a + b - 1) / b }

// cells are the scene's rectangles, shrunk n times.
func (d *drawn) cells(n int) []cell {
	if cells, ok := d.fit[n]; ok {
		return cells
	}
	// Ink far from the background outlasts ink near it.
	var rank [pixel.Slots + 1]int
	back := d.colors[pixel.Back]
	for s, c := range d.colors {
		dr, dg, db := int(c.R)-int(back.R), int(c.G)-int(back.G), int(c.B)-int(back.B)
		rank[s+1] = 1 + dr*dr + dg*dg + db*db
	}
	d.fit[n] = d.sc.shrink(n, &rank).cells()
	slog.Debug("8bit: mermaid fit", "factor", n, "cells", len(d.fit[n]))
	return d.fit[n]
}

// report logs one render, never the source.
func (d *drawn) report() {
	var pe *pele.Error
	switch {
	case d.err == nil:
		slog.Debug("8bit: mermaid drew", "bytes", len(d.source), "width", d.sc.w, "height", d.sc.h)
	case errors.As(d.err, &pe):
		slog.Debug("8bit: mermaid failed", "bytes", len(d.source), "code", pe.Code, "line", pe.Line, "column", pe.Column)
	default:
		slog.Debug("8bit: mermaid failed", "bytes", len(d.source), "code", "walk")
	}
}
