package slice

import (
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	as := []string{"foo", "bar", "bar", "baz"}
	bs := []string{"foo", "bar", "baz"}

	if !reflect.DeepEqual(Unique(as), bs) {
		t.Fatal("didn't remove duplicates")
	}
	if !reflect.DeepEqual(Unique(bs), bs) {
		t.Fatal("didn't preserve list if no duplicates")
	}

	ai := []int{1, 2, 3, 3, 4}
	bi := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(Unique(ai), bi) {
		t.Fatal("didn't remove duplicates")
	}
	if !reflect.DeepEqual(Unique(bi), bi) {
		t.Fatal("didn't preserve list if no duplicates")
	}
}
