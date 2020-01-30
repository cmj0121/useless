/* Copyright (C) 2020-2020 cmj. All right reserved. */
package zargs

import (
	"testing"
)

func TestFlag(t *testing.T) {
	if true != NewFlag("enable", false).Set("").Get().(bool) {
		/* Set the boolean flag failure */
		t.Errorf("Cannot set bool flag")
	}

	if "10" != NewFlag("enable", nil).Set("10").Get().(string) {
		/* Set the string (default) flag failure */
		t.Errorf("Cannot set string (default) flag")
	}

	if "changed" != NewFlag("enable", "default").Set("changed").Get().(string) {
		/* Set the string flag failure */
		t.Errorf("Cannot set string flag")
	}

	if 10 != NewFlag("enable", 1).Set("10").Get().(int) {
		/* Set the integer flag failure */
		t.Errorf("Cannot set int flag")
	}
}

func TestFlagDefault(t *testing.T) {
	if false != NewFlag("enable", false).Get().(bool) {
		/* Set the boolean flag failure */
		t.Errorf("Cannot set default bool flag")
	}

	if "" != NewFlag("enable", nil).Get().(string) {
		/* Set the string (default) flag failure */
		t.Errorf("Cannot set default string (default) flag")
	}

	if "default" != NewFlag("enable", "default").Get().(string) {
		/* Set the string flag failure */
		t.Errorf("Cannot set default string flag")
	}

	if 1 != NewFlag("enable", 1).Get().(int) {
		/* Set the integer flag failure */
		t.Errorf("Cannot set default int flag")
	}
}
