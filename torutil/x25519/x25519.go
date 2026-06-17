// Package x25519 provides X25519 key pair helpers for v3 onion service client
// authorization. The generated keys are used by the "ONION_CLIENT_AUTH_*"
// controller commands (client side, base64-encoded private key) and the
// ADD_ONION "ClientAuthV3" flag / "descriptor:x25519:" auth files (service
// side, base32-encoded public key).
package x25519

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"io"
)

// pubKeyBase32Encoding matches the unpadded, upper-case base32 encoding Tor uses
// for client-auth public keys.
var pubKeyBase32Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// KeyPair is an X25519 key pair used for v3 onion service client authorization.
type KeyPair struct {
	priv *ecdh.PrivateKey
}

// Generate creates a new random X25519 key pair. If rng is nil,
// crypto/rand.Reader is used.
func Generate(rng io.Reader) (*KeyPair, error) {
	if rng == nil {
		rng = rand.Reader
	}
	priv, err := ecdh.X25519().GenerateKey(rng)
	if err != nil {
		return nil, err
	}
	return &KeyPair{priv: priv}, nil
}

// FromPrivateKey builds a KeyPair from a raw 32-byte X25519 private key.
func FromPrivateKey(privateKey []byte) (*KeyPair, error) {
	priv, err := ecdh.X25519().NewPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return &KeyPair{priv: priv}, nil
}

// FromPrivateKeyBase64 builds a KeyPair from a base64-encoded raw X25519 private
// key, the form used by the ONION_CLIENT_AUTH_ADD command.
func FromPrivateKeyBase64(blob string) (*KeyPair, error) {
	byts, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 private key: %w", err)
	}
	return FromPrivateKey(byts)
}

// PrivateKey returns the raw 32-byte X25519 private key.
func (k *KeyPair) PrivateKey() []byte { return k.priv.Bytes() }

// PublicKey returns the raw 32-byte X25519 public key.
func (k *KeyPair) PublicKey() []byte { return k.priv.PublicKey().Bytes() }

// PrivateKeyBase64 returns the base64-encoded private key, the form expected by
// the ONION_CLIENT_AUTH_ADD controller command (client side).
func (k *KeyPair) PrivateKeyBase64() string {
	return base64.StdEncoding.EncodeToString(k.PrivateKey())
}

// PublicKeyBase32 returns the base32-encoded public key, the form expected by
// the ADD_ONION "ClientAuthV3" flag and "descriptor:x25519:" client auth files
// (service side).
func (k *KeyPair) PublicKeyBase32() string {
	return pubKeyBase32Encoding.EncodeToString(k.PublicKey())
}
