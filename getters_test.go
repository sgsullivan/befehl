package befehl

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sgsullivan/befehl/helpers/filesystem"
)

func getZeroValOpts() *Instance {
	if i, e := New(&Options{
		PrivateKeyFile: "",
		LogDir:         "",
		SshHostKeyConfig: SshHostKeyConfig{
			Enabled:        false,
			KnownHostsPath: "",
		},
		RunConfigPath: "unit-test-resources/zero-hosts.json",
	}); e != nil {
		panic(e)
	} else {
		return i
	}
}

func getNonZeroValOpts() *Instance {
	i, err := New(&Options{
		PrivateKeyFile: "foo",
		LogDir:         "baz",
		SshHostKeyConfig: SshHostKeyConfig{
			Enabled:        true,
			KnownHostsPath: defaultKnownHosts,
		},
		RunConfigPath: "unit-test-resources/hosts.json",
	})
	if err != nil {
		panic(err)
	}
	return i
}

var defaultSshPath = os.Getenv("HOME") + "/.ssh"
var defaultKnownHosts = defaultSshPath + "/known_hosts"

func init() {
	if !filesystem.PathExists(defaultSshPath) {
		if err := os.Mkdir(defaultSshPath, os.ModePerm); err != nil {
			panic(fmt.Sprintf("failed to create %s: %s", defaultSshPath, err))
		}
	}

	if !filesystem.FileExists(defaultKnownHosts) {
		f, err := os.Create(defaultKnownHosts)
		if err != nil {
			panic(fmt.Sprintf("failed to create %s: %s", defaultKnownHosts, err))
		}
		f.Close()
	}
}

func TestGetSshUser(t *testing.T) {
	zuser := getZeroValOpts().getDefaultSshUser()
	if zuser != "root" {
		t.Fatalf("User [%s] for zeroval is unexpected", zuser)
	}
	nuser := getNonZeroValOpts().getDefaultSshUser()
	if nuser != "r00t" {
		t.Fatalf("User [%s] for nonzeroval is unexpected", nuser)
	}
}

func TestGetLogDir(t *testing.T) {
	logDirZero := getZeroValOpts().getLogDir()
	if logDirZero != os.Getenv("HOME")+"/befehl/logs" {
		t.Fatalf("LogDir of %s for zeroval is unexpected", logDirZero)
	}
	logDir := getNonZeroValOpts().getLogDir()
	if logDir != "baz" {
		t.Fatalf("LogDir of %s for zeroval is unexpected", logDir)
	}
}

func TestGetLogFilePath(t *testing.T) {
	got := getZeroValOpts().getLogFilePath("server")
	expected := os.Getenv("HOME") + "/befehl/logs" + "/server"
	if got != expected {
		t.Fatalf("getLogFilePath got [%s] expected [%s]", got, expected)
	}
}

func TestGetPrivKeyFile(t *testing.T) {
	if getZeroValOpts().getPrivKeyFile() != os.Getenv("HOME")+"/.ssh/id_rsa" {
		t.Fatal("PrivateKeyFile for zeroval is unexpected")
	}
	if getNonZeroValOpts().getPrivKeyFile() != "foo" {
		t.Fatal("PrivateKeyFile for nonzeroval is unexpected")
	}
}

func TestGetSshClientConfig(t *testing.T) {
	got, err := getNonZeroValOpts().getSshClientConfig(func() string { return getNonZeroValOpts().getDefaultSshUser() })
	if err != nil {
		t.Fatal(err)
	}
	if got.Timeout != time.Duration(10)*time.Second {
		t.Fatalf("returned timeout %s is unexpected", got.Timeout)
	}
}

func TestGetSshKnowHostsPath(t *testing.T) {
	if getZeroValOpts().getSshKnowHostsPath() != os.Getenv("HOME")+"/.ssh/known_hosts" {
		t.Fatal("getSshKnowHostsPath for zeroval is unexpected")
	}
	if getNonZeroValOpts().getSshKnowHostsPath() != defaultKnownHosts {
		t.Fatal("PrivateKeyFile for nonzeroval is unexpected")
	}
}
