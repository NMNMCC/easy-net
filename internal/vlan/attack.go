//go:build linux || android
// +build linux android

package vlan

import (
	"fmt"

	"nmnm.cc/easy-net/internal/log"
)

type AttackConfig struct {
	Start uint16
	Link  string
}

var vlanAttackLogger = log.New("vlan/attack")

func Attack(cfg *AttackConfig) error {
	var start uint16 = cfg.Start
	var count uint16 = 0

	vlanAttackLogger.Info("starting VLAN attack", "link", cfg.Link, "start", start)
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
			vlanAttackLogger.Info("found valid VLAN ID", "id", id)
			return nil
		} else {
			vlanAttackLogger.Warn("failed to test VLAN ID", "id", id, "error", err)
		}

		count++
	}

	vlanAttackLogger.Info("VLAN attack finished, no valid ID found")
	return fmt.Errorf("attack failed")
}
