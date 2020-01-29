/* Copyright (C) 2020-2020 cmj. All right reserved. */
package main

import (
	"fmt"

	"github.com/cmj0121/useless/dcard"
)

func main() {
	agent := dcard.New()

	fmt.Println(agent.Boards()[0])
	fmt.Println(agent.Posts("ntu", -1))
}
