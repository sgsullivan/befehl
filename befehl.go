package befehl

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/fatih/color"
	"github.com/thoas/go-funk"

	"github.com/sgsullivan/befehl/helpers/filesystem"
	"github.com/sgsullivan/befehl/helpers/slice"
	"github.com/sgsullivan/befehl/helpers/waitgroup"
	"github.com/sgsullivan/befehl/queue"
)

func New(options *Options) (*Instance, error) {
	runtimeConfig, err := GetRuntimeConfig(options.RunConfigPath)
	if err != nil {
		return nil, err
	}

	instance := &Instance{
		options:       options,
		runtimeConfig: &runtimeConfig,
	}

	if err := instance.ensureUniqueHosts(); err != nil {
		return nil, err
	}

	return instance, nil
}

func (instance *Instance) Execute(routines int) error {
	if instance.sshKey == nil {
		if err := instance.populateSshKey(); err != nil {
			return err
		}
	}

	return instance.executePayloadOnHosts(routines)
}

func (instance *Instance) ensureUniqueHosts() error {
	maybeDupes := funk.Map(
		instance.runtimeConfig.Hosts,
		func(c RuntimeConfigHost) string {
			return fmt.Sprintf("%s:%d", c.Host, c.Port)
		},
	).([]string)

	noDupes := slice.Unique(maybeDupes)

	if len(maybeDupes) != len(noDupes) {
		return fmt.Errorf("duplicate Host:Port entries in provided configuration")
	}

	return nil
}

func (instance *Instance) executePayloadOnHosts(routines int) error {
	var wg sync.WaitGroup
	hostCnt := len(instance.runtimeConfig.Hosts)
	wg.Add(hostCnt)
	hostsChan := make(chan int, routines)
	queueInstance := new(queue.Queue).New(int64(hostCnt))

	defaultSshConfig, err := instance.getDefaultSshClientConfig()
	if err != nil {
		return err
	}

	for _, hostEntry := range instance.runtimeConfig.Hosts {
		hostsChan <- 1

		hostEntry := hostEntry

		chosenPayloadPath := instance.runtimeConfig.Payload
		if hostEntry.Payload != "" {
			chosenPayloadPath = hostEntry.Payload
		}
		chosenPayload, err := filesystem.ReadFile(chosenPayloadPath)
		if err != nil {
			return err
		}

		sshConfig := defaultSshConfig
		if hostEntry.User != "" {
			sshConfig, err = instance.getSshUserClientConfig(hostEntry.User)
			if err != nil {
				return err
			}
		}

		go func(hostEntry *RuntimeConfigHost) {
			instance.runPayload(&wg, hostEntry.Host, hostEntry.Port, chosenPayload, sshConfig)
			<-hostsChan
			remaining := queueInstance.DecrementCounter()
			color.Magenta(fmt.Sprintf("Remaining: %d / %d\n", remaining, hostCnt))
		}(&hostEntry)
	}

	if waitgroup.WgTimeout(&wg, time.Duration(1800)*time.Second) {
		return fmt.Errorf("hit timeout waiting for all routines to finish")
	}

	color.Green("All routines completed!\n")
	return nil
}

func (instance *Instance) runPayload(wg *sync.WaitGroup, host string, port int, payload []byte, sshConfig *ssh.ClientConfig) {
	defer wg.Done()
	log.Printf("running payload on %s:%d ..\n", host, port)

	hostPort := fmt.Sprintf("%s:%d", host, port)

	// establish the connection
	conn, err := ssh.Dial("tcp", hostPort, sshConfig)
	if err != nil {
		uhoh := fmt.Sprintf("ssh.Dial() to %s failed: %s\n", host, err)
		color.Red(uhoh)
		if err := instance.logPayloadRun(hostPort, uhoh); err != nil {
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
		if err := instance.logPayloadRun(hostPort, uhoh); err != nil {
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
		if err := instance.logPayloadRun(hostPort, uhoh); err != nil {
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
	if err := instance.logPayloadRun(hostPort, cmdOutput); err != nil {
		panic(err)
	}
}
