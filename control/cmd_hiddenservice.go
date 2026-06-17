package control

import "strings"

// GetHiddenServiceDescriptorAsync invokes HSFETCH.
func (c *Conn) GetHiddenServiceDescriptorAsync(address string, server string) error {
	cmd := "HSFETCH " + address
	if server != "" {
		cmd += " SERVER=" + server
	}
	return c.sendRequestIgnoreResponse("%s", cmd)
}

// PostHiddenServiceDescriptorAsync invokes HSPOST.
func (c *Conn) PostHiddenServiceDescriptorAsync(desc string, servers []string, address string) error {
	var cmd strings.Builder
	cmd.WriteString("+HSPOST")
	for _, server := range servers {
		cmd.WriteString(" SERVER=" + server)
	}
	if address != "" {
		cmd.WriteString("HSADDRESS=" + address)
	}
	cmd.WriteString("\r\n" + desc + "\r\n.")
	return c.sendRequestIgnoreResponse("%s", cmd.String())
}
