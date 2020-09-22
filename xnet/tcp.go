// Copyright (c) 2020, Geert JM Vanderkelen

package xnet

import (
	"fmt"
	"net"
)

// GetTCPPort gets a free TCP port for given address. When address
// is not provided (is empty), 127.0.0.1 is used.
func GetTCPPort(address string) (int, error) {
	if address == "" {
		address = "127.0.0.1"
	}
	l, err := net.Listen("tcp", address+":0")
	if err != nil {
		return 0, err
	}
	defer func() { _ = l.Close() }()

	return l.Addr().(*net.TCPAddr).Port, nil
}

// MustGetLocalhostTCPPort gets a free TCP port for IP 127.0.0.1.
// Panics when an error occurs.
func MustGetLocalhostTCPPort() int {
	p, err := GetTCPPort("127.0.0.1")
	if err != nil {
		panic(fmt.Sprintf("TCPPort: failed getting TCP port: %s", err))
	}
	return p
}
