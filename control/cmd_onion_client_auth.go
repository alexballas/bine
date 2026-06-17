package control

import (
	"strings"

	"github.com/alexballas/bine/torutil"
)

// OnionClientAuthKeyType is the key type used for v3 onion service client
// authorization. Only x25519 is currently defined.
const OnionClientAuthKeyType = "x25519"

// OnionClientAuthAddRequest is a set of request params for
// ONION_CLIENT_AUTH_ADD.
type OnionClientAuthAddRequest struct {
	// Address is the v3 onion service address to register a credential for. The
	// ".onion" suffix is optional and stripped if present.
	Address string
	// PrivateKey is the base64-encoded x25519 private key used to authenticate
	// to the service. See torutil/x25519 for helpers to generate one.
	PrivateKey string
	// ClientName is an optional, human-readable nickname for this credential.
	ClientName string
	// Permanent, if true, stores the credential on disk (in ClientOnionAuthDir)
	// so it survives Tor restarts. Otherwise it only lasts for the lifetime of
	// the current Tor process.
	Permanent bool
}

// AddOnionClientAuth invokes ONION_CLIENT_AUTH_ADD to register a client-side
// credential for connecting to a v3 onion service that requires authorization.
func (c *Conn) AddOnionClientAuth(req *OnionClientAuthAddRequest) error {
	if req.PrivateKey == "" {
		return c.protoErr("PrivateKey required")
	}
	var cmd strings.Builder
	cmd.WriteString("ONION_CLIENT_AUTH_ADD ")
	cmd.WriteString(trimOnionSuffix(req.Address))
	cmd.WriteString(" " + OnionClientAuthKeyType + ":" + req.PrivateKey)
	if req.ClientName != "" {
		cmd.WriteString(" ClientName=" + torutil.EscapeSimpleQuotedStringIfNeeded(req.ClientName))
	}
	if req.Permanent {
		cmd.WriteString(" Flags=Permanent")
	}
	resp, err := c.SendRequest("%s", cmd.String())
	// Tor replies 252 when the credential let it decrypt an already-cached
	// descriptor; that is success, but SendRequest surfaces it as an error.
	if err != nil && resp != nil && resp.Err.Code == StatusOkOnionClientAuthDecrypted {
		return nil
	}
	return err
}

// RemoveOnionClientAuth invokes ONION_CLIENT_AUTH_REMOVE to delete a previously
// registered client credential for the given v3 onion address. The ".onion"
// suffix is optional.
func (c *Conn) RemoveOnionClientAuth(address string) error {
	return c.sendRequestIgnoreResponse("ONION_CLIENT_AUTH_REMOVE %s", trimOnionSuffix(address))
}

// OnionClientAuth is a single client credential returned by ViewOnionClientAuth.
type OnionClientAuth struct {
	// Address is the v3 onion service address (without the ".onion" suffix).
	Address string
	// KeyType is the credential key type (e.g. "x25519").
	KeyType string
	// PrivateKey is the base64-encoded private key blob.
	PrivateKey string
	// ClientName is the optional nickname, or empty if none was set.
	ClientName string
	// Flags are any flags set on the credential (e.g. "Permanent").
	Flags []string
}

// ViewOnionClientAuth invokes ONION_CLIENT_AUTH_VIEW. If address is empty, all
// stored credentials are returned; otherwise only the one for that address (if
// any) is returned. The ".onion" suffix on address is optional.
func (c *Conn) ViewOnionClientAuth(address string) ([]*OnionClientAuth, error) {
	cmd := "ONION_CLIENT_AUTH_VIEW"
	if address != "" {
		cmd += " " + trimOnionSuffix(address)
	}
	resp, err := c.SendRequest("%s", cmd)
	if err != nil {
		return nil, err
	}
	var ret []*OnionClientAuth
	for _, data := range resp.Data {
		key, val, _ := torutil.PartitionString(data, ' ')
		if key != "CLIENT" {
			continue
		}
		ret = append(ret, parseOnionClientAuth(val))
	}
	return ret, nil
}

// parseOnionClientAuth parses a CLIENT line body of the form
// "Address KeyType:PrivateKeyBlob [ClientName=Name] [Flags=Flags]". ClientName
// may be a quoted string containing spaces.
func parseOnionClientAuth(val string) *OnionClientAuth {
	auth := &OnionClientAuth{}
	for i, field := range splitQuotedFields(val) {
		switch i {
		case 0:
			auth.Address = field
		case 1:
			auth.KeyType, auth.PrivateKey, _ = torutil.PartitionString(field, ':')
		default:
			k, v, _ := torutil.PartitionString(field, '=')
			switch k {
			case "ClientName":
				auth.ClientName, _ = torutil.UnescapeSimpleQuotedStringIfNeeded(v)
			case "Flags":
				auth.Flags = strings.Split(v, ",")
			}
		}
	}
	return auth
}

// splitQuotedFields splits s on spaces, keeping double-quoted regions (which may
// contain spaces) together as a single field. Backslash escapes within a quoted
// region are honored so an escaped quote does not end the region.
func splitQuotedFields(s string) []string {
	var fields []string
	var cur strings.Builder
	var inQuote, escaping, started bool
	flush := func() {
		if started {
			fields = append(fields, cur.String())
			cur.Reset()
			started = false
		}
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case escaping:
			cur.WriteByte(c)
			escaping = false
		case c == '\\' && inQuote:
			cur.WriteByte(c)
			escaping = true
		case c == '"':
			inQuote = !inQuote
			cur.WriteByte(c)
			started = true
		case c == ' ' && !inQuote:
			flush()
		default:
			cur.WriteByte(c)
			started = true
		}
	}
	flush()
	return fields
}

func trimOnionSuffix(address string) string {
	return strings.TrimSuffix(address, ".onion")
}
