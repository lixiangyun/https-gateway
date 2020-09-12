package util

import (
	"io/ioutil"
)

func SaveToFile(name string, body []byte) error {
	return ioutil.WriteFile(name, body, 0664)
}

func CopyToFile(toName string, fromName string) error {
	body, err := ioutil.ReadFile(fromName)
	if err != nil {
		return err
	}
	return SaveToFile(toName, body)
}