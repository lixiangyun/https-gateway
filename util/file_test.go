package util

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestSaveToFile(t *testing.T){
	for i:=0; i<100; i++ {
		body := GetToken(Rand(10241024))
		err := SaveToFile("test", []byte(body))
		if err != nil {
			t.Error(err.Error())
			break
		}
		body2, err := ioutil.ReadFile("test")
		if err != nil {
			t.Error(err.Error())
			break
		}
		if bytes.Compare([]byte(body), body2) != 0 {
			t.Errorf("diff: %d %d\n", len(body), len(body2))
			break
		}
	}
	os.Remove("test")
}