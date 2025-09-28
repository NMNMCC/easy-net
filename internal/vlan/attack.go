//go:build linux || android
// +build linux android

package vlan

import (
	"fmt"
	"log/slog"
)

type AttackConfig struct {
	Start uint16
	Link  string
}

func Attack(cfg *AttackConfig) error {
	logger := slog.With("component", "vlan", "start", cfg.Start)

	var start uint16 = cfg.Start
	var count uint16 = 0

	logger.Info("starting VLAN attack", "link", cfg.Link, "start", start)
	for {
		if count >= 4095 {
			break
		}

		var id uint16
		if count%2 == 0 {
			id = start + count/2
		} else {
			id = start - (count+1)/2
		}

		if err := Test(&TestConfig{
			Link: cfg.Link,
			ID:   id,
		}); err == nil {
			logger.Info("found valid VLAN ID", "id", id)
			return nil
		}

		count++
	}

	logger.Info("VLAN attack finished, no valid ID found")
	return fmt.Errorf("attack failed")
}
