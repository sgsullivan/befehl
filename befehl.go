package befehl

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"

	"github.com/sgsullivan/befehl/helpers/filesystem"
	"github.com/sgsullivan/befehl/helpers/waitgroup"
)

type queue struct {
	count int64
}

func (q *queue) signifyComplete(total int) {
	remaining := atomic.AddInt64(&q.count, -1)
	color.Magenta(fmt.Sprintf("Remaining: %d / %d\n", remaining, total))
}

type Options struct {
	PrivateKeyFile string
	SshUser        string
	LogDir         string
}

type Instance struct {
	options *Options
	sshKey  ssh.Signer
}

func New(options *Options) *Instance {
	return &Instance{
		options: options,
	}
}

func (instance *Instance) Fire(hostsFile, payload string, routines int) error {
	if bytePayload, readFileErr := filesystem.ReadFile(payload); readFileErr == nil {
		if instance.sshKey != nil {
			if err := instance.populateSshKey(); err != nil {
				return err
			}
		}
		return instance.executePayloadOnHosts(bytePayload, hostsFile, routines)
	} else {
		return readFileErr
	}
}

func (instance *Instance) getPrivKeyFile() string {
	privKeyFile := os.Getenv("HOME") + "/.ssh/id_rsa"
	if instance.options.PrivateKeyFile != "" {
		privKeyFile = instance.options.PrivateKeyFile
	}

	return privKeyFile
}

func (instance *Instance) populateSshKeyEncrypted(privKeyBytes *pem.Block) error {
	fmt.Printf("enter private key password: ")
	password, err := gopass.GetPasswd()
	if err != nil {
		return fmt.Errorf("error when reading input: %v", err)
	}

	pwBuf, err := x509.DecryptPEMBlock(privKeyBytes, []byte(password))
	if err != nil {
		return fmt.Errorf("x509.DecryptPEMBlock failed: %v", err)
	}

	pk, err := x509.ParsePKCS1PrivateKey(pwBuf)
	if err != nil {
		return fmt.Errorf("x509.ParsePKCS1PrivateKey failed: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(pk)
	if err != nil {
		return fmt.Errorf("ssh.NewSignerFromKey failed: %v", err)
	}

	instance.sshKey = signer

	return nil
}

func (instance *Instance) populateSshKeyUnencrypted(rawKey []byte) error {
	signer, err := ssh.ParsePrivateKey(rawKey)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %v", err)
	}

	instance.sshKey = signer

	return nil
}

func (instance *Instance) populateSshKey() error {
	privKeyFile := instance.getPrivKeyFile()

	if rawKey, readFileError := filesystem.ReadFile(privKeyFile); readFileError == nil {
		privKeyBytes, _ := pem.Decode(rawKey)

		if x509.IsEncryptedPEMBlock(privKeyBytes) {
			return instance.populateSshKeyEncrypted(privKeyBytes)
		} else {
			return instance.populateSshKeyUnencrypted(rawKey)
		}
	} else {
		return readFileError
	}
}

func (instance *Instance) getSshUser() string {
	if instance.options.SshUser != "" {
		return instance.options.SshUser
	}
	return "root"
}

func (instance *Instance) buildHostLists(hostsFilePath string) ([]string, error) {
	hostsFile, err := os.Open(hostsFilePath)
	if err != nil {
		return nil, err
	}
	defer hostsFile.Close()

	victims := []string{}

	scanner := bufio.NewScanner(hostsFile)
	for scanner.Scan() {
		host := scanner.Text()
		victims = append(victims, host)
	}

	return victims, scanner.Err()
}

func (instance *Instance) executePayloadOnHosts(payload []byte, hostsFilePath string, routines int) error {
	hostsList, err := instance.buildHostLists(hostsFilePath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	hostCnt := len(hostsList)
	wg.Add(hostCnt)
	var sem = make(chan int, routines)

	sshEntryUser := instance.getSshUser()

	sshConfig := &ssh.ClientConfig{
		User: sshEntryUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(instance.sshKey),
		},
		Timeout:         time.Duration(10) * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	queue := new(queue)
	queue.count = int64(hostCnt)
	for _, host := range hostsList {
		host := host
		sem <- 1
		go func() {
			instance.runPayload(&wg, host, payload, sshConfig)
			<-sem
			queue.signifyComplete(hostCnt)
		}()
	}

	if waitgroup.WgTimeout(&wg, time.Duration(1800)*time.Second) {
		return fmt.Errorf("hit timeout waiting for all routines to finish")
	}

	color.Green("All routines completed!\n")
	return nil
}

func (instance *Instance) runPayload(wg *sync.WaitGroup, host string, payload []byte, sshConfig *ssh.ClientConfig) {
	defer wg.Done()
	log.Printf("running payload on %s ..\n", host)

	// establish the connection
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), sshConfig)
	if err != nil {
		uhoh := fmt.Sprintf("ssh.Dial() to %s failed: %s\n", host, err)
		color.Red(uhoh)
		if err := instance.logPayloadRun(host, uhoh); err != nil {
			panic(err)
		}
		return
	}
	defer conn.Close()

	// open the session
	session, err := conn.NewSession()
	if err != nil {
		uhoh := fmt.Sprintf("ssh.NewSession() to %s failed: %s\n", host, err)
		color.Red(uhoh)
		if err := instance.logPayloadRun(host, uhoh); err != nil {
			panic(err)
		}
		return
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 24, 80, modes); err != nil {
		uhoh := fmt.Sprintf("session.RequestPty() to %s failed: %s\n", host, err)
		color.Red(uhoh)
		if err := instance.logPayloadRun(host, uhoh); err != nil {
			panic(err)
		}
		return
	}

	// finally, run the payload
	var sessionRunAttempt string
	if err := session.Run(string(payload)); err != nil {
		sessionRunAttempt = fmt.Sprintf("session.Run() to %s raised error: %s\n", host, err)
		color.Red(sessionRunAttempt)
	}

	cmdOutput := stdout.String() + stderr.String() + "\n" + sessionRunAttempt
	if err := instance.logPayloadRun(host, cmdOutput); err != nil {
		panic(err)
	}
}

func (instance *Instance) logPayloadRun(host string, output string) error {
	logDir := os.Getenv("HOME") + "/befehl/logs"
	if instance.options.LogDir != "" {
		logDir = instance.options.LogDir
	}
	logFile := logDir + "/" + host
	if !filesystem.PathExists(logDir) {
		if err := os.MkdirAll(logDir, os.FileMode(0700)); err != nil {
			return fmt.Errorf("failed creating [%s]: %s", logDir, err)
		}
	}
	f, err := os.Create(logFile)
	if err != nil {
		return fmt.Errorf("error creating [%s]: %s", logFile, err)
	}
	defer f.Close()

	if _, err = f.WriteString(output); err != nil {
		return fmt.Errorf("error writing to [%s]: %s", logFile, err)
	}

	log.Printf("payload completed on %s! logfile at: %s\n", host, logFile)
	return nil
}
