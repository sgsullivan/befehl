package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cast"

	"github.com/sgsullivan/befehl"
	"github.com/sgsullivan/befehl/helpers/filesystem"
	"github.com/sgsullivan/befehl/integration_tests/util/cmd"
)

func buildApplication(cmdCancel chan bool, workDir string) error {
	if stdout, stderr, err := cmd.RunCmd("make", []string{"build-only"}, cmdCancel, workDir, []string{}); err != nil {
		return fmt.Errorf("failed building application: %s %s %s", stdout.String(), stderr.String(), err)
	}

	return nil
}

func runPayload(cmdCancel chan bool, workDir string, evars []string) error {
	if stdout, stderr, err := cmd.RunCmd(
		"./_exe/befehl",
		[]string{
			"execute",
			"--runconfig",
			"integration_tests/examples/hosts.json",
			"--routines",
			"10",
		},
		cmdCancel,
		workDir,
		evars,
	); err != nil {
		return fmt.Errorf("Failed running befehl execute: %s %s %s", stdout.String(), stderr.String(), err)
	} else {
		strStdout := stdout.String()
		if strings.Contains(strStdout, "failed") {
			return fmt.Errorf("a failure was detected for one or more hosts: %s", strStdout)
		}
	}

	return nil
}

func clearSshdHosts(cmdCancel chan bool, workDir string) error {
	if stdout, stderr, err := cmd.RunCmd("make", []string{"integration-nuke-sshd-hosts"}, cmdCancel, workDir, []string{}); err != nil {
		return fmt.Errorf("failed clearing sshd hosts: %s %s %s", stdout.String(), stderr.String(), err)
	}

	return nil
}

func startSshdHosts(cmdCancel chan bool, workDir string) error {
	if stdout, stderr, err := cmd.RunCmd("make", []string{"integration-start-sshd-hosts"}, cmdCancel, workDir, []string{}); err != nil {
		return fmt.Errorf("failed starting sshd hosts: %s %s %s", stdout.String(), stderr.String(), err)
	}

	return nil
}

func getLogDir() string {
	return os.Getenv("HOME") + "/befehl/logs"
}

func getHostsFilePath(workDir string) string {
	return workDir + "/integration_tests/examples/hosts.json"
}

func verifyWrittenPayloadLogs(workDir string, f func(c befehl.RuntimeConfig) error) error {
	runtimeConfig, err := befehl.GetRuntimeConfig(getHostsFilePath(workDir))
	if err != nil {
		return err
	}

	return f(runtimeConfig)
}

func verifyPayloadLogsPresent(workDir string) error {
	return verifyWrittenPayloadLogs(
		workDir,
		func(runtimeConfig befehl.RuntimeConfig) error {
			logDir := getLogDir()
			for _, hostEntry := range runtimeConfig.Hosts {
				expectedLogPath := fmt.Sprintf("%s/%s:%d", logDir, hostEntry.Host, hostEntry.Port)
				if !filesystem.PathExists(expectedLogPath) {
					return fmt.Errorf("for host %+v expected path %s to exist", hostEntry, expectedLogPath)
				}
			}
			return nil
		},
	)
}

func verifyExpectedPayloadsRan(workDir string) error {
	return verifyWrittenPayloadLogs(
		workDir,
		func(runtimeConfig befehl.RuntimeConfig) error {
			logDir := getLogDir()
			for _, hostEntry := range runtimeConfig.Hosts {
				expectedLogPath := fmt.Sprintf("%s/%s:%d", logDir, hostEntry.Host, hostEntry.Port)
				fileContents, err := ioutil.ReadFile(expectedLogPath)
				if err != nil {
					return fmt.Errorf("failed to open %s: %s", expectedLogPath, err)
				}
				strFileContents := cast.ToString(fileContents)
				if hostEntry.Payload != "" {
					if !strings.Contains(strFileContents, "overrode payload") {
						return fmt.Errorf("payload log at %s doesn't contain overrode payload output", expectedLogPath)
					}
				} else {
					if !strings.Contains(strFileContents, "Hello, world") {
						return fmt.Errorf("payload log at %s doesn't contain non overrode payload output", expectedLogPath)
					}
				}
			}
			return nil
		},
	)
}

func TestIntegration(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	workDir := pwd + "/.."
	privateKeyPath := workDir + "/integration_tests/docker/ssh/id_rsa"
	evars := []string{
		fmt.Sprintf("BEFEHL_SSH_PRIVATEKEYFILE=%s", privateKeyPath),
	}

	cmdCancel := make(chan bool)

	if err := buildApplication(cmdCancel, workDir); err != nil {
		t.Fatal(err)
	}

	if err := clearSshdHosts(cmdCancel, workDir); err != nil {
		t.Fatal(err)
	}

	if err := startSshdHosts(cmdCancel, workDir); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Duration(3) * time.Second)

	defer func() {
		if err := clearSshdHosts(cmdCancel, workDir); err != nil {
			panic(err)
		}
	}()

	if err := runPayload(cmdCancel, workDir, evars); err != nil {
		t.Fatal(err)
	}

	defer func() {
		logDir := getLogDir()
		if err := os.RemoveAll(logDir); err != nil {
			panic(fmt.Sprintf("error deleting %s: %s", logDir, err))
		}
	}()

	if err := verifyPayloadLogsPresent(workDir); err != nil {
		t.Fatal(err)
	}

	if err := verifyExpectedPayloadsRan(workDir); err != nil {
		t.Fatal(err)
	}
}
