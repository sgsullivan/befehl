package befehl

import (
	"golang.org/x/crypto/ssh"
)

type SshHostKeyConfig struct {
	Enabled        bool
	KnownHostsPath string
}

type Options struct {
	PrivateKeyFile   string
	LogDir           string
	SshHostKeyConfig SshHostKeyConfig
	RunConfigPath    string
}

type Instance struct {
	options       *Options
	sshKey        ssh.Signer
	runtimeConfig *RuntimeConfig
}

type RuntimeConfig struct {
	Payload string              `json:"payload"`
	Hosts   []RuntimeConfigHost `json:"hosts"`
	User    string              `json:"user"`
}

type RuntimeConfigHost struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	User    string `json:"user"`
	Payload string `json:"payload"`
}
