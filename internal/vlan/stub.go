//go:build !linux && !android
// +build !linux,!android

package vlan

import "fmt"

type AttackConfig struct {
	Start uint16
	Link  string
}

type TestConfig struct {
	Link string
	ID   uint16
}

func Attack(cfg *AttackConfig) error {
	return fmt.Errorf("vlan attack is only supported on linux and android")
}

func Test(cfg *TestConfig) error {
	return fmt.Errorf("vlan test is only supported on linux and android")
}
