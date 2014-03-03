// Package lz4 implements compression using lz4.c
//
// Copyright (c) 2013 CloudFlare, Inc.

package lz4

// #cgo CFLAGS: -O3
// #include "src/lz4.h"
// #include "src/lz4.c"
import "C"

import (
	"fmt"
	"unsafe"
)

// p gets a char pointer to the first byte of a []byte slice
func p(in []byte) *C.char {
	if len(in) == 0 {
		return (*C.char)(nil)
	}
	return (*C.char)(unsafe.Pointer(&in[0]))
}

// clen gets the length of a []byte slice as a char *
func clen(s []byte) C.int {
	return C.int(len(s))
}

// Uncompress with a known output size. len(out) should be equal to
// the length of the uncompressed outout.
func Uncompress(in []byte, out []byte) (err error) {
	read := int(C.LZ4_uncompress(p(in), p(out), clen(out)))

	if read != len(in) {
		err = fmt.Errorf("Uncompress read %d bytes should have read %d",
			read, len(in))
	}
	return
}

// CompressBound calculates the size of the output buffer needed by
// Compress. This is based on the following macro:
//
// #define LZ4_COMPRESSBOUND(isize)
//      ((unsigned int)(isize) > (unsigned int)LZ4_MAX_INPUT_SIZE ? 0 : (isize) + ((isize)/255) + 16)
func CompressBound(in []byte) int {
	return len(in) + ((len(in) / 255) + 16)
}

// Compress compresses in and puts the content in out. len(out)
// should have enough space for the compressed data (use CompressBound
// to calculate). Returns the number of bytes in the out slice.
func Compress(in []byte, out []byte) (outSize int, err error) {
	outSize = int(C.LZ4_compress_limitedOutput(p(in), p(out), clen(in), clen(out)))
	if outSize == 0 {
		err = fmt.Errorf("Insufficient space for compression")
	}
	return
}