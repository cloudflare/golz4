// Package lz4 implements compression using lz4.c. This is its test
// suite.
//
// Copyright (c) 2013 CloudFlare, Inc.

package lz4

import (
	"io/ioutil"
	"strings"
	"testing"
	"testing/quick"
)

func TestCompressionHCRatio(t *testing.T) {
	input, err := ioutil.ReadFile("sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	output := make([]byte, CompressBound(input))
	outSize, err := CompressHC(input, output)
	if err != nil {
		t.Fatal(err)
	}

	if want := 4317; want != outSize {
		t.Fatalf("HC Compressed output length != expected: %d != %d", want, outSize)
	}
}

func TestCompressionHCLevels(t *testing.T) {
	input, err := ioutil.ReadFile("sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Level   int
		Outsize int
	}{
		{0, 4317},
		{1, 4349},
		{2, 4349},
		{3, 4333},
		{4, 4321},
		{5, 4317},
		{6, 4317},
		{7, 4317},
		{8, 4317},
		{9, 4317},
		{10, 4316},
		{11, 4316},
		{12, 4313},
		{13, 4313},
		{14, 4313},
		{15, 4313},
		{16, 4313},
	}

	for _, tt := range cases {
		output := make([]byte, CompressBound(input))
		outSize, err := CompressHCLevel(input, output, tt.Level)
		if err != nil {
			t.Fatal(err)
		}

		if want := tt.Outsize; want != outSize {
			t.Errorf("HC level %d length != expected: %d != %d",
				tt.Level, want, outSize)
		}
	}
}

func TestCompressionHC(t *testing.T) {
	input := []byte(strings.Repeat("Hello world, this is quite something", 10))
	output := make([]byte, CompressBound(input))
	outSize, err := CompressHC(input, output)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}
	if outSize == 0 {
		t.Fatal("Output buffer is empty.")
	}
	output = output[:outSize]
	decompressed := make([]byte, len(input))
	var ulen int
	ulen, err = Uncompress(output, decompressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}
	if string(decompressed) != string(input) {
		t.Fatalf("Decompressed output != input: %q != %q", decompressed, input)
	}
	if ulen != len(input) {
		t.Fatalf("uncompressed lenght != input length, %v != %v", ulen, len(input))
	}
}

func TestEmptyCompressionHC(t *testing.T) {
	input := []byte("")
	output := make([]byte, CompressBound(input))

	outSize, err := CompressHC(input, output)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}
	if outSize == 0 {
		t.Fatal("Output buffer is empty.")
	}
	output = output[:outSize]
	decompressed := make([]byte, len(input))
	var ulen int
	ulen, err = Uncompress(output, decompressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}
	if string(decompressed) != string(input) {
		t.Fatalf("Decompressed output != input: %q != %q", decompressed, input)
	}
	if ulen != len(input) {
		t.Fatalf("uncompressed lenght != input length, %v != %v", ulen, len(input))
	}
}

func TestNoCompressionHC(t *testing.T) {
	input := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	output := make([]byte, CompressBound(input))
	outSize, err := CompressHC(input, output)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}
	if outSize == 0 {
		t.Fatal("Output buffer is empty.")
	}
	output = output[:outSize]
	decompressed := make([]byte, len(input))
	var ulen int
	ulen, err = Uncompress(output, decompressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}
	if string(decompressed) != string(input) {
		t.Fatalf("Decompressed output != input: %q != %q", decompressed, input)
	}
	if ulen != len(input) {
		t.Fatalf("uncompressed lenght != input length, %v != %v", ulen, len(input))
	}
}

func TestCompressionErrorHC(t *testing.T) {
	input := []byte(strings.Repeat("Hello world, this is quite something", 10))
	output := make([]byte, 0)
	outSize, err := CompressHC(input, output)

	if outSize != 0 {
		t.Fatalf("%d", outSize)
	}

	if err == nil {
		t.Fatalf("Compression should have failed but didn't")
	}

	output = make([]byte, 1)
	_, err = CompressHC(input, output)
	if err == nil {
		t.Fatalf("Compression should have failed but didn't")
	}
}

func TestDecompressionErrorHC(t *testing.T) {
	input := []byte(strings.Repeat("Hello world, this is quite something", 10))
	output := make([]byte, CompressBound(input))
	outSize, err := CompressHC(input, output)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}
	if outSize == 0 {
		t.Fatal("Output buffer is empty.")
	}
	output = output[:outSize]
	decompressed := make([]byte, len(input)-1)
	_, err = Uncompress(output, decompressed)
	if err == nil {
		t.Fatalf("Decompression should have failed")
	}

	decompressed = make([]byte, 1)
	_, err = Uncompress(output, decompressed)
	if err == nil {
		t.Fatalf("Decompression should have failed")
	}

	decompressed = make([]byte, 0)
	_, err = Uncompress(output, decompressed)
	if err == nil {
		t.Fatalf("Decompression should have failed")
	}
}

func TestFuzzHC(t *testing.T) {
	f := func(input []byte) bool {
		output := make([]byte, CompressBound(input))
		outSize, err := CompressHC(input, output)
		if err != nil {
			t.Fatalf("Compression failed: %v", err)
		}
		if outSize == 0 {
			t.Fatal("Output buffer is empty.")
		}
		output = output[:outSize]
		decompressed := make([]byte, len(input))
		var ulen int
		ulen, err = Uncompress(output, decompressed)
		if err != nil {
			t.Fatalf("Decompression failed: %v", err)
		}
		if string(decompressed) != string(input) {
			t.Fatalf("Decompressed output != input: %q != %q", decompressed, input)
		}
		if ulen != len(input) {
			t.Fatalf("uncompressed lenght != input length, %v != %v", ulen, len(input))
		}

		return true
	}

	conf := &quick.Config{MaxCount: 20000}
	if testing.Short() {
		conf.MaxCount = 1000
	}
	if err := quick.Check(f, conf); err != nil {
		t.Fatal(err)
	}
}
