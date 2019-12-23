/* Copyright (C) 2019-2019 cmj. All right reserved. */
package pf

import (
	"os"
)

/* The NOP packer */
type NOPPacker struct {
	data []byte
}

func (p *NOPPacker) Read(in []byte) {
	p.data = in
	return
}

func (p *NOPPacker) Compress() (err error) {
	return
}

func (p *NOPPacker) Obfuscate() (err error) {
	return
}

func (p *NOPPacker) Save(dst string) (err error) {
	var f *os.File

	if f, err = os.Create(dst); err != nil {
		/* Cannot save to destination */
		return
	}

	f.Write(p.data)
	return
}
