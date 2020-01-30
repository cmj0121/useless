/* Copyright (C) 2020-2020 cmj. All right reserved. */
package zargs

import (
	"fmt"
)

type Zargs struct {
	Description string                /* The description shows in the command line */
	Flags       []*Option             /* Extra flags */
	Subs        map[string]*Zargs     /* Sub-commands */
	Callback    func(in *Zargs) error /* Callback function for match the command */

	/* index of the flags */
	idx_flags      map[rune]*Option
	idx_long_flags map[string]*Option
}

type Option interface {
	/* The abstract option in Zargs which may the flag or the argument */
	Get() interface{}
	Set(string)
	String() string
}

func New(description string) (out *Zargs) {
	out = &Zargs{
		Description:    description,
		Subs:           make(map[string]*Zargs, 0),
		idx_flags:      make(map[rune]*Option),
		idx_long_flags: make(map[string]*Option),
	}
	return
}

func (args *Zargs) Parse(in ...string) (err error) {
	idx := 0
	for idx < len(in) {
		switch in[idx][0] {
		case '-':
			if idx, err = args.ParseFlag(in, idx); err != nil {
				/* parse failure */
				return
			}
		default:
			/* parse as the command */
			if idx, err = args.ParseCommand(in, idx); err != nil {
				/* parse failure */
				return
			}
		}
	}
	return
}

func (args *Zargs) ParseFlag(in []string, idx int) (out int, err error) {
	return
}

func (args *Zargs) ParseCommand(in []string, idx int) (out int, err error) {
	return
}

func (args *Zargs) Sub(cmd, description string) (out *Zargs) {
	if _, ok := args.Subs[cmd]; ok == true {
		err := fmt.Errorf("Duplicated sub-command '%s'", cmd)
		panic(err)
	}

	out = New(description)
	args.Subs[cmd] = out
	return
}

func (args *Zargs) Flag(name string, in interface{}) (out *Flag) {
	out = NewFlag(name, in)
	return
}
