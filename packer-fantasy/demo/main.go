/* Copyright (C) 2019-2019 cmj. All right reserved. */
package main

import (
	pf "github.com/cmj0121/useless/packer-fantasy"
	"os"
)

func main() {
	pf.New(os.Args[1], os.Args[2], "nop")
}
