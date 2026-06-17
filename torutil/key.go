package torutil

import (
	"crypto"
	"crypto/sha3"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/alexballas/bine/torutil/ed25519"
)

var serviceIDEncoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// OnionServiceIDFromPrivateKey generates the onion service ID from the given
// private key. This panics if the private key is not a
// github.com/alexballas/bine/torutil/ed25519.KeyPair.
func OnionServiceIDFromPrivateKey(key crypto.PrivateKey) string {
	switch k := key.(type) {
	case ed25519.KeyPair:
		return OnionServiceIDFromV3PublicKey(k.PublicKey())
	}
	panic(fmt.Sprintf("Unrecognized private key type: %T", key))
}

// OnionServiceIDFromPublicKey generates the onion service ID from the given
// public key. This panics if the public key is not a
// github.com/alexballas/bine/torutil/ed25519.PublicKey.
func OnionServiceIDFromPublicKey(key crypto.PublicKey) string {
	switch k := key.(type) {
	case ed25519.PublicKey:
		return OnionServiceIDFromV3PublicKey(k)
	}
	panic(fmt.Sprintf("Unrecognized public key type: %T", key))
}

// OnionServiceIDFromV3PublicKey generates a V3 service ID for the given
// ED25519 public key.
func OnionServiceIDFromV3PublicKey(key ed25519.PublicKey) string {
	checkSum := sha3.Sum256(append(append([]byte(".onion checksum"), key...), 0x03))
	var keyBytes [35]byte
	copy(keyBytes[:], key)
	keyBytes[32] = checkSum[0]
	keyBytes[33] = checkSum[1]
	keyBytes[34] = 0x03
	return strings.ToLower(serviceIDEncoding.EncodeToString(keyBytes[:]))
}

// PublicKeyFromV3OnionServiceID returns a public key for the given service ID
// or an error if the service ID is invalid.
func PublicKeyFromV3OnionServiceID(id string) (ed25519.PublicKey, error) {
	byts, err := serviceIDEncoding.DecodeString(strings.ToUpper(id))
	if err != nil {
		return nil, err
	} else if len(byts) != 35 {
		return nil, fmt.Errorf("Invalid id length")
	} else if byts[34] != 0x03 {
		return nil, fmt.Errorf("Invalid version")
	}
	// Do a checksum check
	key := ed25519.PublicKey(byts[:32])
	checkSum := sha3.Sum256(append(append([]byte(".onion checksum"), key...), 0x03))
	if byts[32] != checkSum[0] || byts[33] != checkSum[1] {
		return nil, fmt.Errorf("Invalid checksum")
	}
	return key, nil
}
