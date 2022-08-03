package filesystem

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

func FileExists(file string) bool {
	_, err := os.Stat(file)

	return !os.IsNotExist(err)
}
