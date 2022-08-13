package slice

import (
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	a := []string{"foo", "bar", "bar", "baz"}
	b := []string{"foo", "bar", "baz"}

	if !reflect.DeepEqual(Unique(a), b) {
		t.Fatal("didn't remove duplicates")
	}
	if !reflect.DeepEqual(Unique(b), b) {
		t.Fatal("didn't preserve list if no duplicates")
	}

}
