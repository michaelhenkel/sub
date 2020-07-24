package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	byteX, err := getProtoSet("desc.protoset")
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, byteX)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	b := fmt.Sprintf("%# 20x", buf.Bytes())
	s := strings.Fields(b)
	f := strings.Join(s, ", ")
	template := `
package main

var DSCByte = []byte{` + f + `}`
	d1 := []byte(template)
	err = ioutil.WriteFile("../../dsc.go", d1, 0644)
	if err != nil {
		panic(err)
	}
}

func getProtoSet(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}
