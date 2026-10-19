// Package files holds browsing, storage and sharing components.
package files

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// File is an entry in a folder.
type File struct {
	Name string
	Kind Kind
	// Size is in bytes.
	Size              int64
	Created, Modified time.Time
}

// size says n bytes in B, KB, MB, GB, TB.
func size(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	v, i := float64(n)/1024, 0
	// From 999.5 KB it reads 1.0 MB, not 1000 KB.
	for v >= 999.5 && i < 4 {
		v, i = v/1024, i+1
	}
	if v < 9.95 {
		return fmt.Sprintf("%.1f %cB", v, "KMGTP"[i])
	}
	return fmt.Sprintf("%.0f %cB", v, "KMGTP"[i])
}

// pick panics unless *i is -1 or points into n.
func pick(name string, i *int, n int) {
	if *i < -1 || *i >= n {
		panic(fmt.Sprintf("files: %s choice %d of %d", name, *i, n))
	}
}

// check panics on a nameless file, bad kind or size.
func check(name string, f File) {
	look(f.Kind)
	if f.Name == "" || f.Size < 0 {
		panic(fmt.Sprintf("files: %s file %q of %d bytes", name, f.Name, f.Size))
	}
}

// row sets a file's icon, name and a note.
func row(c *ui.Context, k pixel.Kit, f File, note func()) {
	ui.Row(c).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
		FileIcon(c, f.Kind, FileIconOptions{})
		pixel.Text(c, f.Name, pixel.TextOptions{Wrap: true, MaxLines: 1}).Grow(1)
		note()
	})
}
