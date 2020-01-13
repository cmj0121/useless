/* Copyright (C) 2019-2020 cmj. All right reserved. */
package tree

import (
	"reflect"
)

type Tree struct {
	*treeView
	In interface{}

	/* Used to avoid the stack-overflow */
	types    map[reflect.Type]struct{}
	colorize bool
}

func New(in interface{}) (out *Tree) {
	out = &Tree{
		In:    in,
		types: make(map[reflect.Type]struct{}, 0),
	}
	return
}

func (in *Tree) Colorize() (out *Tree) {
	in.colorize = true
	out = in
	return
}

func (in *Tree) String() (out string) {
	out = in.ToString(-1)
	return
}

func (in *Tree) ToString(lv int) (out string) {
	/* parse the target object */
	in.treeView = in.Parse(reflect.TypeOf(in.In), lv)
	out = in.treeView.String()
	return
}

func (tree *Tree) Parse(in reflect.Type, lv int) (out *treeView) {
	/* NOTE - update the types map first */
	tree.types[in] = struct{}{}

	switch in.Kind() {
	case reflect.Ptr:
		out = tree.Parse(in.Elem(), lv)
		/* Set as the POINTER */
		node := out.Data.(*nodeType)
		node.Type = in
	case reflect.Struct:
		node := &nodeType{
			FieldName: "",
			Type:      in,
			Colorize:  tree.colorize,
		}
		out = NewTreeView(node)

		/* NOTE - Only recursive when level is not zero */
		if lv != 0 {
			for idx := 0; idx < in.NumField(); idx++ {
				field := in.Field(idx)

				if _, ok := tree.types[field.Type]; ok {
					/* NOTE - avoid stack overflow, skip */
					continue
				}

				field_tree := tree.Parse(field.Type, lv-1)
				out.InsertTree(field_tree)

				if field.Anonymous == false {
					node := field_tree.Data.(*nodeType)
					node.FieldName = field.Name
				}
			}
		}
	default:
		node := &nodeType{
			Type:     in,
			Colorize: tree.colorize,
		}
		out = NewTreeView(node)
	}

	return
}
