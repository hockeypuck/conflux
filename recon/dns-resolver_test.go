package recon

import (
	n "net"
	"strings"
	"testing"
)

func assertResolveHost(t *testing.T, net, addr string) []n.IPAddr {
	res, err := resolveHost(net, addr)
	if err != nil {
		t.Errorf("Error in resolving %s::%s; %s", net, addr, err)
		return nil
	}
	return res
}

func TestResolveHost(t *testing.T) {
	var res []n.IPAddr
	res = assertResolveHost(t, "tcp4", "example.net:11370")
	if len(res) < 1 {
		t.Errorf("Failed to get a result resolving example.net")
	}
	for _, addr := range res {
		if strings.Contains(addr.String(), ":") {
			t.Errorf("Got a non ipv4 address %s", addr)
		}
	}

	res = assertResolveHost(t, "tcp6", "example.net:11370")
	if len(res) < 1 {
		t.Errorf("Failed to get a result resolving example.net")
	}
	for _, addr := range res {
		if strings.Count(addr.String(), ".") != 0 {
			t.Errorf("Got a non ipv6 address %s", addr)
		}
	}
	res = assertResolveHost(t, "tcp", "example.net:11370")
	// we expect a v4 and a v6 address here
	if len(res) < 2 {
		t.Errorf("Failed to get a result resolving example.net")
	}
}

func TestResolveIPv4(t *testing.T) {
	res := assertResolveHost(t, "tcp4", "127.0.0.1:11370")
	if len(res) != 1 {
		t.Errorf("Failed to get a result resolving tcp4/127.0.0.1")
	}
	res = assertResolveHost(t, "tcp6", "127.0.0.1:11370")
	if len(res) != 0 {
		t.Errorf("Failed to get a result resolving tcp6/127.0.0.1")
	}
	res = assertResolveHost(t, "tcp", "127.0.0.1:11370")
	if len(res) != 1 {
		t.Errorf("Failed to get a result resolving tcp/127.0.0.1")
	}
}

func TestResolveIPv6(t *testing.T) {
	res := assertResolveHost(t, "tcp4", "[::1]:11370")
	if len(res) != 0 {
		t.Errorf("Failed to get a result resolving tcp4/127.0.0.1")
	}
	res = assertResolveHost(t, "tcp6", "[::1]:11370")
	if len(res) != 1 {
		t.Errorf("Failed to get a result resolving tcp6/127.0.0.1")
	}
	res = assertResolveHost(t, "tcp", "[::1]:11370")
	if len(res) != 1 {
		t.Errorf("Failed to get a result resolving tcp/127.0.0.1")
	}
}
func TestInvalidNetwork(t *testing.T) {
	res, err := resolveHost("foo", "bar")
	if res != nil {
		t.Errorf("Got an result instead of an eror")
	}
	if err.Error() != "unknown network foo" {
		t.Errorf("Expected a different error: %s", err)
	}
}

func TestMissingPort(t *testing.T) {
	res, err := resolveHost("tcp", "foobar.example.")
	if res != nil {
		t.Errorf("Got an result instead of an eror")
	}
	if err.Error() != "address foobar.example.: missing port in address" {
		t.Errorf("Expected a different error: [%s]", err)
	}
}

func TestInvalidHost(t *testing.T) {
	res, err := resolveHost("tcp", "_foobar_:1234")

	if err == nil {
		t.Error("Expected an error, but got none")
	}
	if strings.Contains(err.Error(), "no such host") != true {
		t.Errorf("Expected a different error: [%s]", err)
	}
	if len(res) != 0 {
		t.Errorf("Unexpected results; %v", res)
	}
}
