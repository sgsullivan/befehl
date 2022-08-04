package befehl

import (
	"reflect"
	"testing"
)

func TestBuildHostList(t *testing.T) {
	res, err := getNonZeroValOpts().buildHostList("unit_test_resources/legacy_hosts")
	if err != nil {
		t.Fatalf("error building hostList: %s", err)
	}
	if len(res) != 9 {
		t.Fatalf("expected host size to be 10 but was %d", len(res))
	}

	expected := []string{
		"192.168.2.2",
		"192.168.2.3",
		"192.168.2.4",
		"192.168.2.5",
		"192.168.2.6",
		"192.168.2.7",
		"192.168.2.8",
		"192.168.2.9",
		"192.168.2.10",
	}

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("got: %+v expected: %+v", res, expected)
	}
}
