package filesystem

import (
	"os"
	"testing"
)

var dummyPath = "foobarbiz"

func TestPathExists(t *testing.T) {
	if testPathExists(dummyPath) {
		t.Fatalf("expected %s to not exist", dummyPath)
	}
	if PathExists(dummyPath) {
		t.Fatalf("returns true for non-existent file")
	}
	if err := os.Mkdir(dummyPath, os.ModePerm); err != nil {
		t.Fatalf("failed creating %s: %s", dummyPath, err)
	}
	defer func() {
		if err := os.Remove(dummyPath); err != nil {
			t.Fatalf("failed removing %s: %s", dummyPath, err)
		}
	}()
	if !PathExists(dummyPath) {
		t.Fatalf("returns false for existent file")
	}
}
