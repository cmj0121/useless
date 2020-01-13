/* Copyright (C) 2019-2020 cmj. All right reserved. */
package tree

import (
	"fmt"
	"reflect"
	"strings"
)

type Tree struct {
	*treeView
	In interface{}

	types map[reflect.Type]struct{}
}

type TreeNode struct {
	Field string
	Type  string
}

func (node TreeNode) String() (out string) {
	if out = fmt.Sprintf("%-8s [%s]", node.Field, node.Type); node.Field == "" {
		out = node.Type
	}
	out = strings.Trim(out, " ")
	return
}

func New(in interface{}) (out *Tree) {
	out = &Tree{
		In:    in,
		types: make(map[reflect.Type]struct{}, 0),
	}
	return
}

func (in *Tree) String() (out string) {
	/* parse the target object */
	in.treeView = in.Parse(reflect.TypeOf(in.In))
	out = in.treeView.String()
	return
}

func (tree *Tree) Parse(in reflect.Type) (out *treeView) {
	tree.types[in] = struct{}{}

	switch in.Kind() {
	case reflect.Ptr:
		out = tree.Parse(in.Elem())
		/* Set as the POINTER */
		node := out.Data.(*TreeNode)
		node.Type = fmt.Sprintf("*%s", node.Type)
	case reflect.Slice, reflect.Array:
		out = tree.Parse(in.Elem())
		/* Set as the slice */
		node := out.Data.(*TreeNode)
		node.Type = fmt.Sprintf("slice %s", node.Type)
	case reflect.Map:
		out = tree.Parse(in.Elem())
		/* Set as the map */
		node := out.Data.(*TreeNode)
		node.Type = fmt.Sprintf("map %s:%s", in.Key(), node.Type)
	case reflect.Chan:
		out = tree.Parse(in.Elem())
		/* Set as the slice */
		node := out.Data.(*TreeNode)
		node.Type = fmt.Sprintf("%s %s", in.ChanDir(), node.Type)
	case reflect.Struct:
		node := &TreeNode{
			Field: "",
			Type:  in.Name(),
		}

		if node.Type == "" {
			/* anonymous structure */
			node.Type = "struct"
		}

		out = NewTreeView(node)
		/* NOTE - update the types map first */

		for idx := 0; idx < in.NumField(); idx++ {
			field := in.Field(idx)

			if _, ok := tree.types[field.Type]; ok {
				/* NOTE - avoid stack overflow */
				continue
			}

			field_tree := tree.Parse(field.Type)
			out.InsertTree(field_tree)

			if field.Anonymous == false {
				node := field_tree.Data.(*TreeNode)
				node.Field = field.Name
			}
		}
	case reflect.Interface:
		node := &TreeNode{
			Field: "",
			Type:  "interface",
		}
		out = NewTreeView(node)
	default:
		node := &TreeNode{
			Field: "",
			Type:  in.Name(),
		}
		out = NewTreeView(node)
	}

	return
}
