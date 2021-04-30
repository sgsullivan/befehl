package befehl

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/spf13/viper"
)

var Config *viper.Viper

type Queue struct {
	count int64
}

var store struct {
	sshKey ssh.Signer
}

func Fire(targets, payload string, routines int, passedViper *viper.Viper) {
	Config = passedViper
	bytePayload := readFile(payload)
	populateSshKey()
	fireTorpedos(bytePayload, targets, routines)
}

func populateSshKey() {
	// do nothing if sshKey Signer is already stored!
	if store.sshKey != nil {
		return
	}

	privKeyFile := os.Getenv("HOME") + "/.ssh/id_rsa"
	if Config.GetString("auth.privatekeyfile") != "" {
		privKeyFile = Config.GetString("auth.privatekeyfile")
	}
	rawKey := readFile(privKeyFile)
	privKeyBytes, _ := pem.Decode(rawKey)

	if x509.IsEncryptedPEMBlock(privKeyBytes) {
		fmt.Printf("enter private key password: ")
		password, err := gopass.GetPasswd()
		if err != nil {
			panic(fmt.Sprintf("Error when reading input: %v", err))
		}
		pwBuf, err := x509.DecryptPEMBlock(privKeyBytes, []byte(password))
		if err != nil {
			panic(fmt.Sprintf("x509.DecryptPEMBlock failed: %v", err))
		}
		pk, err := x509.ParsePKCS1PrivateKey(pwBuf)
		if err != nil {
			panic(fmt.Sprintf("x509.ParsePKCS1PrivateKey failed: %v", err))
		}
		signer, err := ssh.NewSignerFromKey(pk)
		if err != nil {
			panic(fmt.Sprintf("ssh.NewSignerFromKey failed: %v", err))
		}
		store.sshKey = signer
	} else {
		signer, err := ssh.ParsePrivateKey(rawKey)
		if err != nil {
			panic(fmt.Sprintf("unable to parse private key: %v", err))
		}
		store.sshKey = signer
	}
}

func (q *Queue) signifyComplete(total int) {
	remaining := atomic.AddInt64(&q.count, -1)
	color.Magenta(fmt.Sprintf("Remaining: %d / %d\n", remaining, total))
}

func fireTorpedos(payload []byte, targets string, routines int) {
	file, err := os.Open(targets)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hostCnt := 0
	victims := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := scanner.Text()
		victims = append(victims, host)
		hostCnt++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(hostCnt)
	var sem = make(chan int, routines)

	sshEntryUser := "root"
	if Config.GetString("auth.sshuser") != "" {
		sshEntryUser = Config.GetString("auth.sshuser")
	}

	sshConfig := &ssh.ClientConfig{
		User: sshEntryUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(store.sshKey),
		},
		Timeout:         time.Duration(time.Duration(10) * time.Second),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	queue := new(Queue)
	queue.count = int64(hostCnt)
	for _, host := range victims {
		host := host
		sem <- 1
		go func() {
			runPayload(&wg, host, payload, sshConfig)
			<-sem
			queue.signifyComplete(hostCnt)
		}()
	}

	if wgTimeout(&wg, time.Duration(time.Duration(1800)*time.Second)) {
		panic("hit timeout waiting for all routines to finish")
	}
	color.Green("All routines completed!\n")
}

func runPayload(wg *sync.WaitGroup, host string, payload []byte, sshConfig *ssh.ClientConfig) {
	defer wg.Done()
	log.Printf("running payload on %s ..\n", host)

	// establish the connection
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), sshConfig)
	if err != nil {
		uhoh := fmt.Sprintf("ssh.Dial() to %s failed: %s\n", host, err)
		color.Red(uhoh)
		logPayloadRun(host, uhoh)
		return
	}
	defer conn.Close()

	// open the session
	session, err := conn.NewSession()
	if err != nil {
		uhoh := fmt.Sprintf("ssh.NewSession() to %s failed: %s\n", host, err)
		color.Red(uhoh)
		logPayloadRun(host, uhoh)
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
		logPayloadRun(host, uhoh)
		return
	}

	// finally, run the payload
	var sessionRunAttempt string
	if err := session.Run(string(payload)); err != nil {
		sessionRunAttempt = fmt.Sprintf("session.Run() to %s raised error: %s\n", host, err)
		color.Red(sessionRunAttempt)
	}

	cmdOutput := stdout.String() + stderr.String() + "\n" + sessionRunAttempt
	logPayloadRun(host, cmdOutput)
}

func logPayloadRun(host string, output string) {
	logDir := os.Getenv("HOME") + "/befehl/logs"
	if Config.GetString("general.logdir") != "" {
		logDir = Config.GetString("general.logdir")
	}
	logFile := logDir + "/" + host
	if !pathExists(logDir) {
		if err := os.MkdirAll(logDir, os.FileMode(0700)); err != nil {
			panic(fmt.Sprintf("Failed creating [%s]: %s\n", logDir, err))
		}
	}
	f, err := os.Create(logFile)
	if err != nil {
		panic(fmt.Sprintf("Error creating [%s]: %s", logFile, err))
	}
	defer f.Close()

	if _, err = f.WriteString(output); err != nil {
		panic(fmt.Sprintf("Error writing to [%s]: %s", logFile, err))
	}

	log.Printf("payload completed on %s! logfile at: %s\n", host, logFile)
}

func readFile(file string) []byte {
	read, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return read
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func wgTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
