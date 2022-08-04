package befehl

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cast"
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

func (instance *Instance) validateHostEntry(hostEntry string) error {
	if strings.Contains(hostEntry, ":") {
		split := strings.Split(hostEntry, ":")
		if len(split) != 2 {
			return fmt.Errorf("malformed host entry (multiple :): %s", hostEntry)
		}
		if _, err := strconv.Atoi(split[1]); err != nil {
			return fmt.Errorf("malformed host entry (non numeric port): %s", split[1])
		}
	}

	return nil
}

func (instance *Instance) rawSplitHostEntry(hostEntry string) (string, int) {
	split := strings.Split(hostEntry, ":")

	stringPort := "22"
	if len(split) > 1 {
		stringPort = split[1]
	}

	return split[0], cast.ToInt(stringPort)
}

func (instance *Instance) transformHostFromHostEntry(hostEntry string) (host string, port int, err error) {
	if err = instance.validateHostEntry(hostEntry); err != nil {
		return
	}

	host, port = instance.rawSplitHostEntry(hostEntry)

	return
}
