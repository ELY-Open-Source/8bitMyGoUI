package code

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Var is a name, type and value; structs have fields.
type Var struct {
	Name, Type, Value string
	Fields            []Var
}

type VariablesOptions struct {
	// Width is in dots; zero means 200.
	Width    int
	Disabled bool
}

type VariablesResult = data.TreeResult

// Variables lists names, types and values; fields open below.
func Variables(c *ui.Context, text string, vars []Var, selected *string, o VariablesOptions) VariablesResult {
	return data.Tree(c, text, nodes(vars), selected, data.TreeOptions{Width: pixel.Or("Variables Width", o.Width, 200), Disabled: o.Disabled})
}

// nodes are vars as tree nodes, values colored by kind.
func nodes(vars []Var) []data.TreeNode {
	var out []data.TreeNode
	for _, v := range vars {
		out = append(out, data.TreeNode{Label: v.Name + ": " + v.Type, Value: v.Value, Slot: kind(v), Children: nodes(v.Fields)})
	}
	return out
}

// kind is a value's slot, read from its type.
func kind(v Var) pixel.Slot {
	t := v.Type
	switch {
	case v.Value == "nil":
		return pixel.Faint
	case t == "string":
		return pixel.Good
	case t == "bool":
		return pixel.Info
	case t == "byte" || t == "rune" || strings.HasPrefix(t, "int") || strings.HasPrefix(t, "uint") || strings.HasPrefix(t, "float") || strings.HasPrefix(t, "complex"):
		return pixel.Warn
	}
	return pixel.Faint
}
