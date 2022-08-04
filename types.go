package befehl

import (
	"golang.org/x/crypto/ssh"
)

type Options struct {
	PrivateKeyFile string
	SshUser        string
	LogDir         string
}

type Instance struct {
	options *Options
	sshKey  ssh.Signer
}
