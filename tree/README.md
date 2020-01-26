# Tree #
The Go-lang based tree library: show the structure and the nested fields.
The **tree** will show the fields for the golang object and display as the
tree view. For example:

```go
package main

import (
	"fmt"

	"github.com/cmj0121/useless/tree"
)

func main() {
	obj := tree.Tree{}
	fmt.Println(tree.New(&obj))
}
```


## Roadmap ##
- [x] show the tree-view of the object's fields
- [x] colorize
- [x] show the methods
- [x] specified the level of the nested fields (default: -1)
