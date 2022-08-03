package filesystem

import (
	"os"
	"testing"
)

var expectedFileContents = "hello world"

func pathExists(path string) bool {
	if _, e := os.Stat(path); e == nil {
		return true
	}

	return false
}

func TestFileExists(t *testing.T) {
	filePath := "foobar"
	if pathExists(filePath) {
		t.Fatalf("filePath [%s] exists", filePath)
	}
	if e := os.WriteFile(
		filePath,
		[]byte(expectedFileContents),
		0644,
	); e != nil {
		t.Fatalf("failed to write to file [%s]", filePath)
	}
	defer func() {
		if err := os.Remove(filePath); err != nil {
			t.Fatalf("failed removing filePath [%s]: %s", filePath, err)
		}
	}()

	if !FileExists(filePath) {
		t.Fatal("incorrectly says file that exists doesnt")
	}

}

func TestReadFile(t *testing.T) {
	dummyFile := "/foo/bar/baz/biz"
	if pathExists(dummyFile) {
		t.Fatalf("expected [%s] to not exist but it does", dummyFile)
	}
	if _, e := ReadFile(dummyFile); e == nil {
		t.Fatal("ReadFile didn't return an error for non-existent file")
	}

	dummyWrittenFile := "foobarbazbiz"
	if e := os.WriteFile(
		dummyWrittenFile,
		[]byte(expectedFileContents),
		0644,
	); e != nil {
		t.Fatalf("failed to write to file [%s]", dummyWrittenFile)
	}
	defer func() {
		if err := os.Remove(dummyWrittenFile); err != nil {
			t.Fatalf("Failed to remove [%s]: %s", dummyWrittenFile, err)
		}
	}()
	if read, err := ReadFile(dummyWrittenFile); err == nil {
		if string(read) != expectedFileContents {
			t.Fatalf("ReadFile returned unexpected data: %s", read)
		}
	} else {
		t.Fatalf("failed to ReadFile: %s", err)
	}
}
