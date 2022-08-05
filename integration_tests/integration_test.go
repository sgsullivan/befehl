package integration

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

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
			"--hosts",
			"integration_tests/examples/hosts",
			"--payload",
			"integration_tests/examples/payload",
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

func verifyPayloadLogsPresent(workDir string) error {
	logDir := os.Getenv("HOME") + "/befehl/logs"

	hostsFilePath := workDir + "/integration_tests/examples/hosts"
	hostsFile, err := os.Open(hostsFilePath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(hostsFile)
	defer func() {
		if err := os.RemoveAll(logDir); err != nil {
			panic(fmt.Sprintf("error deleting %s: %s", logDir, err))
		}
		hostsFile.Close()
	}()
	for scanner.Scan() {
		hostEntry := scanner.Text()
		expectedLogPath := logDir + "/" + hostEntry
		if !filesystem.PathExists(expectedLogPath) {
			return fmt.Errorf("for host %s expected path %s to exist", hostEntry, expectedLogPath)
		}
	}

	return nil
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

	if err := verifyPayloadLogsPresent(workDir); err != nil {
		t.Fatal(err)
	}
}
