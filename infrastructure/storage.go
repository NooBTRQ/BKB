package infrastructure

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"time"

	"github.com/errors"
)

func WriteFile(data []byte, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "open File error")
	}
	n, err := f.Write(data)

	for err != nil {
		data = data[n:]
		n, err = f.Write(data)
	}
	return nil
}

func AppendFile(data []byte, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "open File error")
	}
	n, err := f.Write(data)

	for err != nil {
		data = data[n:]
		n, err = f.Write(data)
	}
	return nil
}

func ReadFile(filePath string) ([]byte, error) {

	return ioutil.ReadFile(filePath)
}

func IntToBytes(n interface{}) ([]byte, error) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.BigEndian, &n)
	return bytesBuffer.Bytes(), err
}

func BytesToInt(b []byte, n interface{}) error {
	bytesBuffer := bytes.NewBuffer(b)
	return binary.Read(bytesBuffer, binary.BigEndian, n)
}

func RuntimeInfoLog(message, logInfoPaht string) {
	message = time.Now().String() + "  " + message + "\r\n"
	b := []byte(message)
	AppendFile(b, logInfoPaht)
}
