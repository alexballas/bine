package control

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOnionClientAuth(t *testing.T) {
	t.Run("full line", func(t *testing.T) {
		auth := parseOnionClientAuth(`vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd x25519:ZGVhZGJlZWY= ClientName="my client" Flags=Permanent`)
		require.Equal(t, "vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd", auth.Address)
		require.Equal(t, "x25519", auth.KeyType)
		require.Equal(t, "ZGVhZGJlZWY=", auth.PrivateKey)
		require.Equal(t, "my client", auth.ClientName)
		require.Equal(t, []string{"Permanent"}, auth.Flags)
	})

	t.Run("minimal line", func(t *testing.T) {
		auth := parseOnionClientAuth("someaddress x25519:ZGVhZGJlZWY=")
		require.Equal(t, "someaddress", auth.Address)
		require.Equal(t, "x25519", auth.KeyType)
		require.Equal(t, "ZGVhZGJlZWY=", auth.PrivateKey)
		require.Empty(t, auth.ClientName)
		require.Empty(t, auth.Flags)
	})
}

func TestTrimOnionSuffix(t *testing.T) {
	require.Equal(t, "abc", trimOnionSuffix("abc.onion"))
	require.Equal(t, "abc", trimOnionSuffix("abc"))
}
