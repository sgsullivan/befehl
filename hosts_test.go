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

func TestValidateHostEntry(t *testing.T) {
	if err := getNonZeroValOpts().validateHostEntry("192.168.2.2"); err != nil {
		t.Fatalf("error for valid host entry (without port): %s", err)
	}
	if err := getNonZeroValOpts().validateHostEntry("192.168.2.2:1000"); err != nil {
		t.Fatalf("error for valid host entry (with port): %s", err)
	}
	if err := getNonZeroValOpts().validateHostEntry("192.168.2.2:"); err == nil {
		t.Fatal("no error for invalid host entry (missing port)")
	}
	if err := getNonZeroValOpts().validateHostEntry("192.168.2.2:2s34"); err == nil {
		t.Fatal("no error for invalid host entry (non numeric port)")
	}
}

func TestRawSplitHostEntry(t *testing.T) {
	ip := "10.0.0.42"
	host, port := getNonZeroValOpts().rawSplitHostEntry(ip)
	if host != ip {
		t.Fatalf("unexpected host; got: %s expected: %s", host, ip)
	}
	if port != 22 {
		t.Fatalf("unexpected port; got: %d expected: 22", port)
	}

	ipWithPort := "10.0.0.43:1244"
	host, port = getNonZeroValOpts().rawSplitHostEntry(ipWithPort)
	if host != "10.0.0.43" {
		t.Fatalf("unexpected host; got: %s expected: 10.0.0.43", host)
	}
	if port != 1244 {
		t.Fatalf("unexpected port; got: %d expected: 1244", port)
	}
}

func TestTransformHostFromHostEntry(t *testing.T) {
	hostEntryInvalid := "192.168.0.2:12s3"
	_, _, err := getNonZeroValOpts().transformHostFromHostEntry(hostEntryInvalid)
	if err == nil {
		t.Fatal("didn't return an error for an invalid port")
	}

	hostEntry := "192.168.0.2:2222"
	host, port, err := getNonZeroValOpts().transformHostFromHostEntry(hostEntry)
	if err != nil {
		t.Fatalf("error for valid hostEntry %s: %s", hostEntry, err)
	}
	if host != "192.168.0.2" {
		t.Fatalf("returned host is %s not 192.168.0.2", host)
	}
	if port != 2222 {
		t.Fatalf("returned port is %d not 2222", port)
	}

	hostEntryNoPort := "192.168.0.2"
	host, port, err = getNonZeroValOpts().transformHostFromHostEntry(hostEntryNoPort)
	if err != nil {
		t.Fatalf("error for valid hostEntry %s: %s", hostEntry, err)
	}
	if host != "192.168.0.2" {
		t.Fatalf("returned host is %s not 192.168.0.2", host)
	}
	if port != 22 {
		t.Fatalf("returned port is %d not 22", port)
	}

}
