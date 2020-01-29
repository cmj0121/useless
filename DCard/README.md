# DCard API #
This is the Go-lang based DCard API for the fetch and analysis.

## Simple Code ##
```go
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
```

## Roadmap ##
- [x] Fetch the posts
	- [ ] Fetch the full content of the post
	- [x] Fetch the comments of the post
- [x] Fetch the boards
