package code

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type JSONTreeOptions struct {
	// Width is in dots; zero means 200.
	Width    int
	Disabled bool
}

type JSONTreeResult = data.TreeResult

// JSONTree shows JSON as a tree; values colored by type.
func JSONTree(c *ui.Context, text, doc string, selected *string, o JSONTreeOptions) JSONTreeResult {
	dec := json.NewDecoder(strings.NewReader(doc))
	dec.UseNumber()
	root := node(dec, "")
	if _, err := dec.Token(); err != io.EOF {
		panic(fmt.Sprintf("code: JSONTree: after the value, %v", err))
	}
	// An empty container or lone value is its own row.
	roots := root.Children
	if len(roots) == 0 {
		roots = []data.TreeNode{root}
	}
	return data.Tree(c, text, roots, selected, data.TreeOptions{Width: pixel.Or("JSONTree Width", o.Width, 200), Disabled: o.Disabled})
}

// node reads the next value under label; bad JSON panics.
func node(dec *json.Decoder, label string) data.TreeNode {
	tok, err := dec.Token()
	if err != nil {
		panic(fmt.Sprintf("code: JSONTree: %v", err))
	}
	n := data.TreeNode{Label: label, Slot: pixel.Faint}
	switch v := tok.(type) {
	case json.Delim:
		end := "]"
		if v == '{' {
			end = "}"
		}
		for i := 0; dec.More(); i++ {
			key := strconv.Itoa(i)
			if v == '{' {
				t, _ := dec.Token()
				key = t.(string)
			}
			n.Children = append(n.Children, node(dec, key))
		}
		dec.Token()
		n.Value = fmt.Sprintf("%c%d%s", v, len(n.Children), end)
		// An empty object or array still opens nothing.
		if n.Children == nil {
			n.Children = []data.TreeNode{}
		}
	case string:
		n.Value, n.Slot = strconv.QuoteToGraphic(v), pixel.Good
	case json.Number:
		n.Value, n.Slot = v.String(), pixel.Warn
	case bool:
		n.Value, n.Slot = strconv.FormatBool(v), pixel.Info
	case nil:
		n.Value = "null"
	}
	return n
}
