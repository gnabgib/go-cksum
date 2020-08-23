// Package cksum copyright 2020, gnabgib
// Use of this source code is governed by the LICENSE file
//
// Tableless version of the cksum algorithm, smaller binaries, probably slower execution
// - Note the type (crcNt) and the helpers (streamNt, bytesNt) are not exported,
//   you're not expected to use this implementation but rhather cksum.go
package cksum

import (
	"bufio"
	"io"
)

const poly = 0x04C11DB7

// crc32Nt - A cksum CRC32 storage structure without tables
type crc32Nt struct {
	p, r  uint32
	Size  int
	final bool
}

// newCrcNt - Create a new CRC
func newCrcNt(poly uint32) *crc32Nt {
	return &crc32Nt{poly, 0, 0, false}
}

// Add a new byte to the running CRC
func (pr *crc32Nt) Add(b byte) {
	if pr.final {
		return
	}
	for i := 7; i >= 0; i-- {
		msb := pr.r & (1 << 31)
		pr.r = pr.r << 1
		if msb != 0 {
			pr.r = pr.r ^ pr.p
		}
	}
	pr.r ^= uint32(b)
	pr.Size++
}

// Check - get the final code back, note once this is called you cannot further `Add`
func (pr *crc32Nt) Check() uint32 {
	if pr.final {
		return pr.r
	}
	for m := pr.Size; ; {
		pr.Add(byte(m) & 0xff)
		pr.Size-- //Don't count this towards size
		m = m >> 8
		if m == 0 {
			break
		}
	}
	pr.Add(0)
	pr.Add(0)
	pr.Add(0)
	pr.Add(0)
	pr.final = true //Prevent further modification
	pr.Size -= 4
	pr.r = ^pr.r
	return pr.r
}

// streamNt - Calculate checksum of a stream without using a table (unexported)
func streamNt(r *bufio.Reader) (uint32, int, error) {
	pr := newCrcNt(poly)
	for done := false; !done; {
		switch b, err := r.ReadByte(); err {
		case io.EOF:
			done = true
		case nil:
			pr.Add(b)
		default:
			return 0, 0, err
		}
	}
	return pr.Check(), pr.Size, nil
}

// bytesNt - Calculate checksum of a byte slice without using a table (unexported)
func bytesNt(data []byte) (uint32, int, error) {
	pr := newCrcNt(poly)
	for _, b := range data {
		pr.Add(b)
	}
	return pr.Check(), pr.Size, nil
}
