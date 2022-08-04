package befehl

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/fatih/color"

	"github.com/sgsullivan/befehl/helpers/filesystem"
	"github.com/sgsullivan/befehl/helpers/waitgroup"
	"github.com/sgsullivan/befehl/queue"
)

func New(options *Options) *Instance {
	return &Instance{
		options: options,
	}
}

func (instance *Instance) Execute(hostsFile, payload string, routines int) error {
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

func (instance *Instance) executePayloadOnHosts(payload []byte, hostsFilePath string, routines int) error {
	hostsList, err := instance.buildHostList(hostsFilePath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	hostCnt := len(hostsList)
	wg.Add(hostCnt)
	hostsChan := make(chan int, routines)
	queueInstance := new(queue.Queue).New(int64(hostCnt))

	sshConfig := instance.getSshClientConfig()

	for _, host := range hostsList {
		host := host
		hostsChan <- 1
		go func() {
			instance.runPayload(&wg, host, payload, sshConfig)
			<-hostsChan
			remaining := queueInstance.DecrementCounter()
			color.Magenta(fmt.Sprintf("Remaining: %d / %d\n", remaining, hostCnt))
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
