/* Copyright (C) 2020-2020 cmj. All right reserved. */
package zargs

import (
	"fmt"
	"reflect"
)

type Flag struct {
	reflect.Type

	Name    string        /* Flag name */
	Default reflect.Value /* default value */

	help     string
	shortcut rune /* The shortcut of the flag */
	required bool /* Required flags */
	raw      string
	value    reflect.Value
}

func NewFlag(name string, in interface{}) (out *Flag) {
	out = &Flag{
		Name: name,
	}

	if in == nil {
		in = ""
	}

	switch v := reflect.ValueOf(in); v.Kind() {
	case reflect.Bool, reflect.Int, reflect.String:
		out.Type = v.Type()
	default:
		err := fmt.Errorf("Not support '%s' (%s)", in, v.Kind())
		panic(err)
	}

	return
}

func (flag *Flag) String() (out string) {
	out = fmt.Sprintf("--%s", flag.Name)
	return
}

func (flag *Flag) Set(in string) {
	flag.raw = in
	return
}

func (flag *Flag) Get() (out interface{}) {
	out = flag.value.Interface()
	return
}

func (flag *Flag) Shortcut(c rune) (out *Flag) {
	flag.shortcut = c
	out = flag
	return
}

func (flag *Flag) Help(in string) (out *Flag) {
	flag.help = in
	out = flag
	return
}
