package system

import (
	"io/ioutil"
	"os"
)

func ReadFile(file string) ([]byte, error) {
	read, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return read, nil
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
