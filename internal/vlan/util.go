//go:build linux || android
// +build linux android

package vlan

import (
	"fmt"
)

func NewVlanName(link string, id uint16) string {
	return fmt.Sprintf("%.10s.%d", link, id)
}
