package cksum

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func expect(t *testing.T, name string, expected, found interface{}) bool {
	if expected == found {
		return true
	}
	t.Errorf("%s was incorrect, expected %v, got %v", name, expected, found)
	return false
}

func bytesTest(t *testing.T, input []byte, expectSize int, expectCrc uint32) {
	crc, size, err := Bytes(input)
	expect(t, "Err", nil, err)
	expect(t, "Crc", expectCrc, crc)
	expect(t, "Size", expectSize, size)
}

func fileTest(t *testing.T, filename string, expectSize int, expectCrc uint32) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Unable to find file", filename)
		return
	}
	defer f.Close()
	in := bufio.NewReader(f)
	crc, size, err := Stream(in)
	expect(t, "Err", nil, err)
	expect(t, "Crc", uint32(2330645186), crc)
	//Arguable whether we need to test the found size
	expect(t, "Size", 4, size)
}

func TestBytes(t *testing.T) {
	bytesTest(t, []byte("123\n"), 4, 2330645186)
	bytesTest(t, []byte(""), 0, 4294967295)
	bytesTest(t, []byte("\n"), 1, 3515105045)
	bytesTest(t, []byte("CRC helps with bit rot\n"), 23, 3193580682)
	bytesTest(t, []byte("I do not want to work\n"), 22, 17471322)
}

func TestFile(t *testing.T) {
	fileTest(t, "eg-file.txt", 4, 2330645186)
}

/*If you wish to extend tests:
run `echo $DATA > test && cksum test`
- Where $DATA is whatever you'd like the content of the string/file to be
- Note that on linux a \n will be appended to your $DATA
- Record the results (sum, size, file and update these tests
*/
