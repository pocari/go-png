package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func dumpChunk(r io.Reader) {
	var length int32
	binary.Read(r, binary.BigEndian, &length)
	bytes := make([]byte, 4)
	r.Read(bytes)

	fmt.Printf("chunk '%v' (%d bytes)\n", string(bytes), length)

}

func readChunks(file *os.File) []io.Reader {
	var chunks []io.Reader

	var offset int64 = 8
	file.Seek(offset, 0)
	for {
		var len int32
		err := binary.Read(file, binary.BigEndian, &len)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		chunks = append(chunks, io.NewSectionReader(file, offset, int64(len)+12))
		offset, err = file.Seek(int64(len+8), 1)
		if err != nil {
			panic(err)
		}
	}
	return chunks
}

func main() {
	f, err := os.Open("./lenna.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	chunks := readChunks(f)
	for _, chunk := range chunks {
		dumpChunk(chunk)
	}
}
