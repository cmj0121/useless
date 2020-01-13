/* Copyright (C) 2019-2020 cmj. All right reserved. */
package tree

import (
	"fmt"
	"reflect"
)

const (
	_RESET   = "\x1b[m"
	_BLACK   = "\x1b[1;30m"
	_RED     = "\x1b[1;31m"
	_GREEN   = "\x1b[1;32m"
	_YELLOW  = "\x1b[1;33m"
	_BLUE    = "\x1b[1;34m"
	_MAGENTA = "\x1b[1;35m"
	_CYAN    = "\x1b[1;36m"
	_WHITE   = "\x1b[1;37m"
	_GRAY    = "\x1b[37m"
)

type nodeType struct {
	reflect.Type

	FieldName string
	Colorize  bool
}

func (node *nodeType) String() (out string) {
	if out = node.TypeName(node.Type); node.FieldName != "" {
		out = fmt.Sprintf("%s%-8s%s [%s]", node.FieldColor(node.FieldName), node.FieldName, node.Reset(), out)
	}
	return
}

func (node *nodeType) TypeName(in reflect.Type) (out string) {
	switch in.Kind() {
	case reflect.Ptr:
		out = fmt.Sprintf("%s*%s%s", node.PointerColor(), node.TypeName(in.Elem()), node.Reset())
	case reflect.Slice:
		out = fmt.Sprintf("%sslice%s %s", node.TagColor(), node.Reset(), node.TypeName(in.Elem()))
	case reflect.Map:
		out = fmt.Sprintf("%smap%s %s:%s", node.TagColor(), node.Reset(), in.Key(), node.TypeName(in.Elem()))
	case reflect.Chan:
		out = fmt.Sprintf("%s %s", in.ChanDir(), node.TypeName(in.Elem()))
	case reflect.Interface:
		out = fmt.Sprintf("%sinterface%s", node.TagColor(), node.Reset())
	case reflect.Struct:
		if out = in.Name(); out == "" {
			/* anonymous structure */
			out = fmt.Sprintf("%sstruct%s", node.TagColor(), node.Reset())
		}
	default:
		out = fmt.Sprintf("%s%s%s", node.TagColor(), in.Name(), node.Reset())
	}

	return
}

func (node *nodeType) Reset() (out string) {
	if out = ""; node.Colorize {
		out = _RESET
	}
	return
}

func (node *nodeType) TagColor() (out string) {
	if out = ""; node.Colorize {
		out = _CYAN
	}
	return
}

func (node *nodeType) FieldColor(in string) (out string) {
	if out = ""; node.Colorize {
		if out = _GRAY; in[0] >= 'A' && in[0] <= 'Z' {
			out = _BLUE
		}
	}
	return
}

func (node *nodeType) PointerColor() (out string) {
	if out = ""; node.Colorize {
		out = _RED
	}
	return
}
