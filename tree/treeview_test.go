/* Copyright (C) 2019-2020 cmj. All right reserved. */
package tree

import (
	"testing"
)

func TestTreeView(t *testing.T) {
	view_1 := NewTreeView("sample")
	if "sample" != view_1.String() {
		/* Cannot show the single tree view */
		t.Errorf("Tree View failure on single node")
	}

	return
}
