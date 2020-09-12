package util

import (
	"io/ioutil"
	"strconv"
	"strings"
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

func LoadPidFile(pid string) int {
	body, err := ioutil.ReadFile(pid)
	if err != nil {
		return 0
	}
	body2 := strings.ReplaceAll(string(body), "\r","")
	body2 = strings.ReplaceAll(body2, "\n","")
	num, err := strconv.Atoi(body2)
	if err != nil {
		return 0
	}
	return num
}