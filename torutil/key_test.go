package torutil

import (
	"encoding/base64"
	"testing"

	"github.com/alexballas/bine/torutil/ed25519"
	"github.com/stretchr/testify/require"
)

func genEd25519(t *testing.T) ed25519.KeyPair {
	k, e := ed25519.GenerateKey(nil)
	require.NoError(t, e)
	return k
}

func TestOnionServiceIDFromPrivateKey(t *testing.T) {
	assert := func(key any, shouldPanic bool) {
		if shouldPanic {
			require.Panics(t, func() { OnionServiceIDFromPrivateKey(key) })
		} else {
			require.NotPanics(t, func() { OnionServiceIDFromPrivateKey(key) })
		}
	}
	assert(nil, true)
	assert("bad type", true)
	assert(genEd25519(t), false)
}

func TestOnionServiceIDFromPublicKey(t *testing.T) {
	assert := func(key any, shouldPanic bool) {
		if shouldPanic {
			require.Panics(t, func() { OnionServiceIDFromPublicKey(key) })
		} else {
			require.NotPanics(t, func() { OnionServiceIDFromPublicKey(key) })
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
		require.NoError(t, err)
		pubKey := ed25519.PrivateKey(key).PublicKey()
		matchingID := matchingIDs[i]
		require.Equal(t, matchingID, OnionServiceIDFromV3PublicKey(pubKey))
		// Check verify here too
		derivedPubKey, err := PublicKeyFromV3OnionServiceID(matchingID)
		require.NoError(t, err)
		require.Equal(t, pubKey, derivedPubKey)
		// Let's mangle the matchingID a bit
		tooLong := matchingID + "ddddd"
		_, err = PublicKeyFromV3OnionServiceID(tooLong)
		require.EqualError(t, err, "Invalid id length")
		badVersion := matchingID[:len(matchingID)-1] + "e"
		_, err = PublicKeyFromV3OnionServiceID(badVersion)
		require.EqualError(t, err, "Invalid version")
		badChecksum := []byte(matchingID)
		badChecksum[len(badChecksum)-3] = 'q'
		_, err = PublicKeyFromV3OnionServiceID(string(badChecksum))
		require.EqualError(t, err, "Invalid checksum")
	}
}
