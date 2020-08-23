# Posix cksum implementation in go

The [CRC-32 checksum](https://en.wikipedia.org/wiki/Cyclic_redundancy_check), and the checksum generated for [Posix cksum](https://en.wikipedia.org/wiki/Cksum) are slightly different despite starting with the same polynomial.  Maybe this is why people people have moved on to [MD5](https://en.wikipedia.org/wiki/MD5) or [SHA](https://en.wikipedia.org/wiki/SHA-1) based checksums.

If you wish to generated *CRC-32 checksums* the [hash/CRC32](https://golang.org/pkg/hash/crc32/) package (using the [ChecksumIEEE](https://golang.org/src/hash/crc32/crc32.go?s=7544:7581#L241) method) will work just fine.  You can compare results of this with the online tool found [here](http://zorc.breitbandkatze.de/crc.html)

For a *posix cksum* version, I've instead implemented the [C code found here](https://pubs.opengroup.org/onlinepubs/009695399/utilities/cksum.html) in golang (a table based calculator).  For reference I also built a non-table version (see: cksum-nt.go)

## Usage

```go
package main

import (
    "fmt"
    "github.com/gnabgib/go-cksum"
    "os"
)

func main() {
    file := "eg-file.txt"
    f,err := os.Open(file)
    if err != nil {
		fmt.Println("Unable to find file",filename)
		return
    }
    defer f.Close()
	in := bufio.NewReader(f)
    crc, size, err := Stream(in)
    if err!=nil {
        fmt.Println("CRC error:",err)
        return
    }
    fmt.Printf("%d %d %s",crc,size,file)
}
```
