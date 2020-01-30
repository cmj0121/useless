/* Copyright (C) 2020-2020 cmj. All right reserved. */
package zargs

import (
	"fmt"
	"reflect"
	"strconv"
)

type Flag struct {
	reflect.Value

	Name    string        /* Flag name */
	Default reflect.Value /* default value */

	help     string
	shortcut rune /* The shortcut of the flag */
	required bool /* Required flags */
	raw      string
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
		out.Value = v
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

func (flag *Flag) Set(in string) (out *Flag) {
	flag.raw = in

	switch flag.Value.Kind() {
	case reflect.Bool:
		flag.Value = reflect.ValueOf(!flag.Value.Bool())
	case reflect.Int:
		v, err := strconv.ParseInt(in, 10, 32)
		if err != nil {
			/* Cannot convert the int */
			err = fmt.Errorf("Set '%s' failure - %s", flag.Name, err)
			panic(err)
		}
		flag.Value = reflect.ValueOf(int(v))
	case reflect.String:
		flag.Value = reflect.ValueOf(in)
	}

	out = flag
	return
}

func (flag *Flag) Get() (out interface{}) {
	out = flag.Value.Interface()
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

func (flag *Flag) Required() (out *Flag) {
	flag.required = true
	out = flag
	return
}
