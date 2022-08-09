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
	SshUser          string
	LogDir           string
	SshHostKeyConfig SshHostKeyConfig
}

type Instance struct {
	options *Options
	sshKey  ssh.Signer
}

type RuntimeConfig struct {
	Payload string              `json:"payload"`
	Hosts   []RuntimeConfigHost `json:"hosts"`
}

type RuntimeConfigHost struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	User    string `json:"user"`
	Payload string `json:"payload"`
}
