package util

import (
	"io/ioutil"
)

func SaveToFile(name string, body []byte) error {
	return ioutil.WriteFile(name, body, 0664)
}