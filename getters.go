package befehl

import (
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func (instance *Instance) getSshUser() string {
	if instance.options.SshUser != "" {
		return instance.options.SshUser
	}
	return "root"
}

func (instance *Instance) getSshClientConfig() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: instance.getSshUser(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(instance.sshKey),
		},
		Timeout:         time.Duration(10) * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func (instance *Instance) getLogDir() string {
	if instance.options.LogDir != "" {
		return instance.options.LogDir
	}
	return os.Getenv("HOME") + "/befehl/logs"
}

func (instance *Instance) getLogFilePath(host string) string {
	return instance.getLogDir() + "/" + host
}

func (instance *Instance) getPrivKeyFile() string {
	if instance.options.PrivateKeyFile != "" {
		return instance.options.PrivateKeyFile
	}

	return os.Getenv("HOME") + "/.ssh/id_rsa"
}
