// Copyright (c) 2020, Geert JM Vanderkelen

package xvenv

import (
	"io/ioutil"
	"strings"
)

func linuxControlGroup() string {
	data, err := ioutil.ReadFile("/proc/1/cgroup")
	if err != nil {
		return ""
	}

	return string(data)
}

// InContainer returns true if the application is running within a container.
func InContainer() bool {
	if strings.Contains(linuxControlGroup(), ":/docker/") {
		return true
	}
	return false
}
