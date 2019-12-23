/* Copyright (C) 2019-2019 cmj. All right reserved. */
package pf

import (
	"fmt"
	"io/ioutil"
)

var binaryFactor map[string]BinaryFantasy
var packerFactor map[string]PackerFantasy

func init() {
	binaryFactor = make(map[string]BinaryFantasy, 0)
	packerFactor = make(map[string]PackerFantasy, 0)

	/* register the default packer */
	RegisterPacker("nop", &NOPPacker{})
}

/* This is used to pack the binary and generate the new packed binary
 *
 *	- Using the BinaryFantasy to detect and analysis the binary
 *	- Using the PackerFantasy to compress and obfuscate machine code
 */
func New(src, dst, packer string) (err error) {
	var data []byte

	if data, err = ioutil.ReadFile(src); err != nil {
		/* Cannot read the source file */
		return
	}

	pf, ok := packerFactor[packer]
	if ok == false {
		/* packer name not found */
		err = fmt.Errorf("Cannot find the packer %s", packer)
		return
	}

	pf.Read(data)
	pf.Compress()
	pf.Obfuscate()
	pf.Save(dst)

	return
}

func RegisterPacker(name string, in PackerFantasy) {
	if _, ok := packerFactor[name]; ok == true {
		panic(fmt.Errorf("Duplicate Packer Fantasy '%s' register", name))
		return
	}

	packerFactor[name] = in
	return
}

func RegisterBinaryFantasy(name string, in BinaryFantasy) {
	if _, ok := binaryFactor[name]; ok == true {
		panic(fmt.Errorf("Duplicate binary Fantasy '%s' register", name))
		return
	}

	binaryFactor[name] = in
	return
}

type BinaryFantasy interface {
	/* Read the file content and raise error when file is not support */
	Read([]byte) error

	/* Write the packed payload to the new machine code */
	Write([]byte)

	/* Save the packed binary to destination path */
	Save(string)

	/* The code that need to packed */
	GetCode() []byte

	/* combine the packed payload with the decoder header */
	Polymorfic(decoder []byte, payload ...[]byte)
}

type PackerFantasy interface {
	/* Read the source content and the BinaryFantasy */
	Read([]byte)

	/* Compress the payload */
	Compress() error

	/* Obfuscate the payload  */
	Obfuscate() error

	/* Save to the destination file path */
	Save(path string) error
}
