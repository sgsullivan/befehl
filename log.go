package befehl

import (
	"fmt"
	"log"
	"os"

	"github.com/sgsullivan/befehl/helpers/filesystem"
)

func (instance *Instance) prepareLogDir() error {
	logDir := instance.getLogDir()
	if !filesystem.PathExists(logDir) {
		if err := os.MkdirAll(logDir, os.FileMode(0700)); err != nil {
			return fmt.Errorf("failed creating [%s]: %s", logDir, err)
		}
	}
	return nil
}

func (instance *Instance) logPayloadRun(host string, output string) error {
	instance.prepareLogDir()
	logFilePath := instance.getLogFilePath(host)

	logFile, err := os.Create(logFilePath)
	if err != nil {
		return fmt.Errorf("error creating [%s]: %s", logFilePath, err)
	}
	defer logFile.Close()

	if _, err = logFile.WriteString(output); err != nil {
		return fmt.Errorf("error writing to [%s]: %s", logFilePath, err)
	}

	log.Printf("payload completed on %s! logfile at: %s\n", host, logFilePath)
	return nil
}
