package control

import (
	"reflect"
	"testing"
)

func TestParseOnionClientAuth(t *testing.T) {
	t.Run("full line", func(t *testing.T) {
		auth := parseOnionClientAuth(`vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd x25519:ZGVhZGJlZWY= ClientName="my client" Flags=Permanent`)
		if auth.Address != "vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd" {
			t.Errorf("Address = %q", auth.Address)
		}
		if auth.KeyType != "x25519" {
			t.Errorf("KeyType = %q, want x25519", auth.KeyType)
		}
		if auth.PrivateKey != "ZGVhZGJlZWY=" {
			t.Errorf("PrivateKey = %q", auth.PrivateKey)
		}
		if auth.ClientName != "my client" {
			t.Errorf("ClientName = %q, want %q", auth.ClientName, "my client")
		}
		if !reflect.DeepEqual(auth.Flags, []string{"Permanent"}) {
			t.Errorf("Flags = %#v, want [Permanent]", auth.Flags)
		}
	})

	t.Run("minimal line", func(t *testing.T) {
		auth := parseOnionClientAuth("someaddress x25519:ZGVhZGJlZWY=")
		if auth.Address != "someaddress" {
			t.Errorf("Address = %q, want someaddress", auth.Address)
		}
		if auth.KeyType != "x25519" {
			t.Errorf("KeyType = %q, want x25519", auth.KeyType)
		}
		if auth.PrivateKey != "ZGVhZGJlZWY=" {
			t.Errorf("PrivateKey = %q", auth.PrivateKey)
		}
		if auth.ClientName != "" {
			t.Errorf("ClientName = %q, want empty", auth.ClientName)
		}
		if len(auth.Flags) != 0 {
			t.Errorf("Flags = %#v, want empty", auth.Flags)
		}
	})
}

func TestTrimOnionSuffix(t *testing.T) {
	if got := trimOnionSuffix("abc.onion"); got != "abc" {
		t.Errorf("trimOnionSuffix(%q) = %q, want abc", "abc.onion", got)
	}
	if got := trimOnionSuffix("abc"); got != "abc" {
		t.Errorf("trimOnionSuffix(%q) = %q, want abc", "abc", got)
	}
}
