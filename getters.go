package befehl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func (instance *Instance) getDefaultSshUser() string {
	if instance.runtimeConfig.User != "" {
		return instance.runtimeConfig.User
	}
	return "root"
}

func (instance *Instance) getSshKnowHostsPath() string {
	path := os.Getenv("HOME") + "/.ssh/known_hosts"

	config := instance.options.SshHostKeyConfig
	if config.Enabled && config.KnownHostsPath != "" {
		path = config.KnownHostsPath
	}

	return path
}

func (instance *Instance) getSshHostKeyCallback() (hostKeyCallback ssh.HostKeyCallback, err error) {
	hostKeyCallback = ssh.InsecureIgnoreHostKey()
	if instance.options.SshHostKeyConfig.Enabled {
		hostKeyCallback, err = knownhosts.New(instance.getSshKnowHostsPath())
	}

	return
}

func (instance *Instance) getSshClientConfig(getSshUser func() string) (*ssh.ClientConfig, error) {
	if hostKeyCallback, err := instance.getSshHostKeyCallback(); err == nil {
		return &ssh.ClientConfig{
			User: getSshUser(),
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(instance.sshKey),
			},
			Timeout:         time.Duration(10) * time.Second,
			HostKeyCallback: hostKeyCallback,
		}, nil
	} else {
		return nil, err
	}
}

func (instance *Instance) getDefaultSshClientConfig() (*ssh.ClientConfig, error) {
	return instance.getSshClientConfig(func() string { return instance.getDefaultSshUser() })
}

func (instance *Instance) getSshUserClientConfig(sshUser string) (*ssh.ClientConfig, error) {
	return instance.getSshClientConfig(func() string { return sshUser })
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

func GetRuntimeConfig(pathToRuntimeConfig string) (config RuntimeConfig, err error) {
	configBytes, err := ioutil.ReadFile(pathToRuntimeConfig)
	if err != nil {
		return
	}

	err = json.Unmarshal(configBytes, &config)

	return
}
