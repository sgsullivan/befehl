package befehl

import (
	"bufio"
	"os"
)

func (instance *Instance) buildHostList(hostsFilePath string) ([]string, error) {
	hostsFile, err := os.Open(hostsFilePath)
	if err != nil {
		return nil, err
	}
	defer hostsFile.Close()

	hostsList := []string{}

	scanner := bufio.NewScanner(hostsFile)
	for scanner.Scan() {
		host := scanner.Text()
		hostsList = append(hostsList, host)
	}

	return hostsList, scanner.Err()
}
