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
