package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

func dumpChunk(r io.Reader) {
	var length int32
	binary.Read(r, binary.BigEndian, &length)
	bytes := make([]byte, 4)
	r.Read(bytes)

	data := make([]byte, length)
	chunkType := string(bytes)

	if chunkType == "tEXt" {
		if _, err := r.Read(data); err != nil {
			panic(err)
		}
	}
	fmt.Printf("chunk '%v' (%d bytes)", chunkType, length)
	if chunkType == "tEXt" {
		fmt.Printf(" data[%s]\n", string(data))
	} else {
		fmt.Println()
	}
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

func textChunk(text string) io.Reader {
	byteData := []byte(text)
	var buffer bytes.Buffer

	hoge := len(byteData)

	binary.Write(&buffer, binary.BigEndian, int32(len(byteData)))
	// これだと動かないので注意
	//binary.Write(&buffer, binary.BigEndian, len(byteData))
	buffer.WriteString("tEXt")
	buffer.Write(byteData)

	crc := crc32.NewIEEE()
	io.WriteString(crc, "tEXt")
	binary.Write(&buffer, binary.BigEndian, crc.Sum32())

	return &buffer
}

func main() {
	f, err := os.Open("./lenna.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	newFile, err := os.Create("./lenna2.png")
	if err != nil {
		panic(err)
	}

	// pngのヘッダ書き込み
	io.WriteString(newFile, "\x89PNG\r\n\x1a\n")

	chunks := readChunks(f)
	io.Copy(newFile, chunks[0])
	io.Copy(newFile, textChunk("Ascii Programming+++"))
	for _, chunk := range chunks[1:] {
		io.Copy(newFile, chunk)
	}
	newFile.Close()

	// 今作ったファイルを読んで見る
	lenna2, err := os.Open("./lenna2.png")
	if err != nil {
		panic(err)
	}
	defer lenna2.Close()

	fmt.Println("newFile chunks")
	newChunks := readChunks(lenna2)
	for _, chunk := range newChunks {
		dumpChunk(chunk)
	}
}
