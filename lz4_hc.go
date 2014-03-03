// Package lz4 implements compression using lz4.c
//
// Copyright (c) 2013 CloudFlare, Inc.

package lz4

// #cgo CFLAGS: -O3
// #include "src/lz4hc.h"
// #include "src/lz4hc.c"
import "C"

import (
	"fmt"
)

func CompressHC(in []byte, out []byte) (outSize int, err error) {
	return CompressHCLevel(in, out, 0)
}

func CompressHCLevel(in []byte, out []byte, level int) (outSize int, err error) {
	// LZ4HC does not handle empty buffers. Pass through to regular Compress.
	if len(in) == 0 || len(out) == 0 {
		return Compress(in, out)
	}

	outSize = int(C.LZ4_compressHC2_limitedOutput(p(in), p(out), clen(in), clen(out), C.int(level)))
	if outSize == 0 {
		err = fmt.Errorf("Insufficient space for compression")
	}
	return
}
