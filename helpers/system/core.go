package system

import (
	"io/ioutil"
	"os"
)

func ReadFileUnsafe(file string) []byte {
	read, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return read
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
