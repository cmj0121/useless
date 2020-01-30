/* Copyright (C) 2020-2020 cmj. All right reserved. */
package main

import (
	"os"

	"github.com/cmj0121/useless/zargs"
)

func main() {
	args := zargs.New("demo program")

	args.Flag("enable", false).Shortcut('e').Help("Enable the flag")

	args.Parse(os.Args[1:]...)
}
