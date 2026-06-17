package x25519

import (
	"encoding/base32"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRoundTrip(t *testing.T) {
	kp, err := Generate(nil)
	require.NoError(t, err)
	require.Len(t, kp.PrivateKey(), 32)
	require.Len(t, kp.PublicKey(), 32)

	// Private key base64 round-trips back to the same key pair.
	kp2, err := FromPrivateKeyBase64(kp.PrivateKeyBase64())
	require.NoError(t, err)
	require.Equal(t, kp.PrivateKey(), kp2.PrivateKey())
	require.Equal(t, kp.PublicKey(), kp2.PublicKey())
}

func TestEncodings(t *testing.T) {
	kp, err := Generate(nil)
	require.NoError(t, err)

	// Base64 private decodes to the raw private key.
	priv, err := base64.StdEncoding.DecodeString(kp.PrivateKeyBase64())
	require.NoError(t, err)
	require.Equal(t, kp.PrivateKey(), priv)

	// Base32 public is unpadded and decodes to the raw public key.
	require.NotContains(t, kp.PublicKeyBase32(), "=")
	pub, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(kp.PublicKeyBase32())
	require.NoError(t, err)
	require.Equal(t, kp.PublicKey(), pub)
}

func TestFromPrivateKeyInvalid(t *testing.T) {
	_, err := FromPrivateKey([]byte("too short"))
	require.Error(t, err)
	_, err = FromPrivateKeyBase64("not!base64!")
	require.Error(t, err)
}
