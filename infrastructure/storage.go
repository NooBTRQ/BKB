package infrastructure

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
)

func WriteFile(data []byte, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	n, err := f.Write(data)

	for err != nil {
		data = data[n:]
		n, err = f.Write(data)
	}
	return err
}

func ReadFile(filePath string) ([]byte, error) {

	return ioutil.ReadFile(filePath)
}

//整形转换成字节
func IntToBytes(n interface{}) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &n)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte, n interface{}) {
	bytesBuffer := bytes.NewBuffer(b)
	binary.Read(bytesBuffer, binary.BigEndian, &n)
}
