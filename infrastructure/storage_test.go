package infrastructure

import "testing"

func TestBytesToInt(t *testing.T) {

	var a int16 = 15
	res, err := IntToBytes(a)

	if err != nil {
		t.Errorf(err.Error())
	}
	println(res)
}
