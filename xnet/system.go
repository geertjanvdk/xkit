// Copyright (c) 2020, Geert JM Vanderkelen

package xnet

import (
	"net"
)

// OutboundIP retrieves the IP address with which the system is communicating
// with the internet.
func OutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:443")
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	return conn.LocalAddr().(*net.UDPAddr).IP, nil
}
