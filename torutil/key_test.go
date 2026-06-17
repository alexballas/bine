package torutil

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/alexballas/bine/torutil/ed25519"
)

func genEd25519(t *testing.T) ed25519.KeyPair {
	t.Helper()
	k, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// didPanic reports whether calling f resulted in a panic.
func didPanic(f func()) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	f()
	return
}

func TestOnionServiceIDFromPrivateKey(t *testing.T) {
	assert := func(key any, shouldPanic bool) {
		t.Helper()
		if got := didPanic(func() { OnionServiceIDFromPrivateKey(key) }); got != shouldPanic {
			t.Errorf("OnionServiceIDFromPrivateKey(%T): panicked=%v, want %v", key, got, shouldPanic)
		}
	}
	assert(nil, true)
	assert("bad type", true)
	assert(genEd25519(t), false)
}

func TestOnionServiceIDFromPublicKey(t *testing.T) {
	assert := func(key any, shouldPanic bool) {
		t.Helper()
		if got := didPanic(func() { OnionServiceIDFromPublicKey(key) }); got != shouldPanic {
			t.Errorf("OnionServiceIDFromPublicKey(%T): panicked=%v, want %v", key, got, shouldPanic)
		}
	}
	assert(nil, true)
	assert("bad type", true)
	assert(genEd25519(t), true)
	assert(genEd25519(t).Public(), false)
}

func TestOnionServiceIDFromV3PublicKey(t *testing.T) {
	base64Keys := []string{
		"SLne6D/uawqUj23619GbeYCd6HnzYPqyUvF8/xyz/3XNVpkgnonQI+J5NQVSGkppD1b0M87+qOtUBmVXsd7H3w",
		"kPUs5aPoqISZVbg0q7coW+mNCODlcL4O7k2QWFOCC0gOQBiDm+g4Xz48lqucA7o2HIQ3gBdL5rlB6+q1tFdJwQ",
		"YGzw/EwpcqfWb5UWIw652Ps4vTKu38VgX7Qo16XvOWjNWQK9YmfgARYiGQ1XYXEAKBJvoq8x+rKFbQN3FG1F6w",
		"IJIZcWE57n5WCvHU2x7GkpBCIw0S0vWd+QyrE5RifGwPtYsbtxjyOxlb754Z0zXLZc+yQUp9hMQt5dt/YNpMag",
		"SD7d4I6ZOjNlcqR2g4ptFJUw0tUHPQvfk92sExvnJ1uofPw9T9LUaaEs3rE/1yoGWKI4YejAzaTJXF9wrWQyuA",
	}
	matchingIDs := []string{
		"2s2wk473fmotzgh6l2ycigrwegnurlzufatjm3bglrb36zbvlerskxad",
		"tmcpdbgklpbywqyjpr7fijvjl7qjihd7pyubosbeohefec2m2thvzoqd",
		"nrcan5uye2fwazixubug6pzrzp6ofjez43bjcyfoxhgxyygxbhgs4zqd",
		"g2csv4kavhunvs45vxxc5ljz775d5a4ycqo4m4nrwpk3b4gryvz2zdyd",
		"jviaiibaz7r6wqxttj5i2bi4zjfilmsevplwwtxdfyjph2sdmq5osdid",
	}
	for i, base64Key := range base64Keys {
		key, err := base64.RawStdEncoding.DecodeString(base64Key)
		if err != nil {
			t.Fatalf("decode key %d: %v", i, err)
		}
		pubKey := ed25519.PrivateKey(key).PublicKey()
		matchingID := matchingIDs[i]
		if got := OnionServiceIDFromV3PublicKey(pubKey); got != matchingID {
			t.Errorf("OnionServiceIDFromV3PublicKey = %q, want %q", got, matchingID)
		}
		// Check verify here too
		derivedPubKey, err := PublicKeyFromV3OnionServiceID(matchingID)
		if err != nil {
			t.Fatalf("PublicKeyFromV3OnionServiceID(%q): %v", matchingID, err)
		}
		if !bytes.Equal(pubKey, derivedPubKey) {
			t.Errorf("derived public key does not match original for id %q", matchingID)
		}
		// Let's mangle the matchingID a bit
		tooLong := matchingID + "ddddd"
		if _, err = PublicKeyFromV3OnionServiceID(tooLong); err == nil || err.Error() != "Invalid id length" {
			t.Errorf("PublicKeyFromV3OnionServiceID(%q) error = %v, want \"Invalid id length\"", tooLong, err)
		}
		badVersion := matchingID[:len(matchingID)-1] + "e"
		if _, err = PublicKeyFromV3OnionServiceID(badVersion); err == nil || err.Error() != "Invalid version" {
			t.Errorf("PublicKeyFromV3OnionServiceID(%q) error = %v, want \"Invalid version\"", badVersion, err)
		}
		badChecksum := []byte(matchingID)
		badChecksum[len(badChecksum)-3] = 'q'
		if _, err = PublicKeyFromV3OnionServiceID(string(badChecksum)); err == nil || err.Error() != "Invalid checksum" {
			t.Errorf("PublicKeyFromV3OnionServiceID(%q) error = %v, want \"Invalid checksum\"", string(badChecksum), err)
		}
	}
}
