package befehl

import (
	"os"
	"testing"
	"time"
)

func getZeroValOpts() *Instance {
	return New(&Options{
		PrivateKeyFile: "",
		SshUser:        "",
		LogDir:         "",
	})
}

func getNonZeroValOpts() *Instance {
	return New(&Options{
		PrivateKeyFile: "foo",
		SshUser:        "bar",
		LogDir:         "baz",
	})
}

func TestGetSshUser(t *testing.T) {
	if getZeroValOpts().getSshUser() != "root" {
		t.Fatal("SshUser for zeroval is unexpected")
	}
	if getNonZeroValOpts().getSshUser() != "bar" {
		t.Fatal("SshUser for nonzeroval is unexpected")
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
	got := getNonZeroValOpts().getSshClientConfig()
	if got.Timeout != time.Duration(10)*time.Second {
		t.Fatalf("returned timeout %s is unexpected", got.Timeout)
	}
}