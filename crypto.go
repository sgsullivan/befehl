package befehl

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/howeyc/gopass"

	"github.com/sgsullivan/befehl/helpers/filesystem"
)

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
