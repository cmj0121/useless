/* Copyright (C) 2020-2020 cmj. All right reserved. */
package zargs

import (
	"testing"
)

func TestFlag(t *testing.T) {
	NewFlag("enable", false)
	NewFlag("enable", nil)
	NewFlag("enable", 10)
	NewFlag("enable", "10")
}
