package x25519

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"strings"
	"testing"
)

func TestGenerateRoundTrip(t *testing.T) {
	kp, err := Generate(nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if l := len(kp.PrivateKey()); l != 32 {
		t.Errorf("PrivateKey length = %d, want 32", l)
	}
	if l := len(kp.PublicKey()); l != 32 {
		t.Errorf("PublicKey length = %d, want 32", l)
	}

	// Private key base64 round-trips back to the same key pair.
	kp2, err := FromPrivateKeyBase64(kp.PrivateKeyBase64())
	if err != nil {
		t.Fatalf("FromPrivateKeyBase64: %v", err)
	}
	if !bytes.Equal(kp.PrivateKey(), kp2.PrivateKey()) {
		t.Errorf("private keys differ after round trip")
	}
	if !bytes.Equal(kp.PublicKey(), kp2.PublicKey()) {
		t.Errorf("public keys differ after round trip")
	}
}

func TestEncodings(t *testing.T) {
	kp, err := Generate(nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// Base64 private decodes to the raw private key.
	priv, err := base64.StdEncoding.DecodeString(kp.PrivateKeyBase64())
	if err != nil {
		t.Fatalf("decode private base64: %v", err)
	}
	if !bytes.Equal(kp.PrivateKey(), priv) {
		t.Errorf("decoded private key does not match raw private key")
	}

	// Base32 public is unpadded and decodes to the raw public key.
	if strings.Contains(kp.PublicKeyBase32(), "=") {
		t.Errorf("PublicKeyBase32 should be unpadded, got %q", kp.PublicKeyBase32())
	}
	pub, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(kp.PublicKeyBase32())
	if err != nil {
		t.Fatalf("decode public base32: %v", err)
	}
	if !bytes.Equal(kp.PublicKey(), pub) {
		t.Errorf("decoded public key does not match raw public key")
	}
}

func TestFromPrivateKeyInvalid(t *testing.T) {
	if _, err := FromPrivateKey([]byte("too short")); err == nil {
		t.Errorf("FromPrivateKey with short key: expected error, got nil")
	}
	if _, err := FromPrivateKeyBase64("not!base64!"); err == nil {
		t.Errorf("FromPrivateKeyBase64 with invalid input: expected error, got nil")
	}
}
