/* Copyright (C) 2019-2020 cmj. All right reserved. */
package main

import (
	"fmt"

	"github.com/cmj0121/useless/tree"
)

func main() {
	obj := tree.Tree{}
	help := tree.New(&obj).Colorize()
	fmt.Println("==== Show to full properties ====")
	fmt.Println(help)

	fmt.Println("==== Show to properties within 1 level ====")
	fmt.Println(help.ToString(1))
}
