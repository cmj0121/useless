/* Copyright (C) 2019-2020 cmj. All right reserved. */
package main

import (
	"fmt"

	"github.com/cmj0121/useless/tree"
)

func main() {
	obj := tree.Tree{}
	help := tree.New(&obj).Colorize()
	fmt.Println(help)
}
