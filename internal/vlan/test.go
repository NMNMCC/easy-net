//go:build linux || android
// +build linux android

package vlan

import (
	"fmt"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/vishvananda/netlink"
	"nmnm.cc/easy-net/internal/log"
)

type TestConfig struct {
	Link    string
	ID      uint16
	Timeout time.Duration
}

var vlanTestLogger = log.New("vlan/test")

func Test(cfg *TestConfig) error {
	master, err := netlink.LinkByName(cfg.Link)
	if err != nil {
		return fmt.Errorf("failed to get link by name: %w", err)
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = NewVlanName(cfg.Link, cfg.ID)
	attrs.ParentIndex = master.Attrs().Index

	vlan := &netlink.Vlan{
		LinkAttrs: attrs,
		VlanId:    int(cfg.ID),
	}

	vlanTestLogger.Info("adding VLAN interface", "link", cfg.Link)
	if err := netlink.LinkAdd(vlan); err != nil {
		return fmt.Errorf("failed to add VLAN interface %s: %w", vlan.Name, err)
	}
	defer func() {
		vlanTestLogger.Info("removing VLAN interface", "name", vlan.Name)
		if err := netlink.LinkDel(vlan); err != nil {
			vlanTestLogger.Error("failed to remove VLAN interface", "name", vlan.Name, "error", err)
			return
		}
		vlanTestLogger.Info("removed VLAN interface", "name", vlan.Name)
	}()
	vlanTestLogger.Info("added VLAN interface", "name", vlan.Name)

	// vlanTestLogger.Info("setting master and bringing up", "name", vlan.Name)
	// if err := netlink.LinkSetMaster(vlan, master); err != nil {
	// 	vlanTestLogger.Error("failed to set master for VLAN interface", "name", vlan.Name, "error", err)
	// 	return err
	// }
	// vlanTestLogger.Info("set master for VLAN interface", "name", vlan.Name)

	vlanTestLogger.Info("bringing up VLAN interface", "name", vlan.Name)
	if err := netlink.LinkSetUp(vlan); err != nil {
		return fmt.Errorf("failed to bring up VLAN interface %s: %w", vlan.Name, err)
	}
	vlanTestLogger.Info("brought up VLAN interface", "name", vlan.Name)

	client := &client4.Client{
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	vlanTestLogger.Info("exchanging DHCP", "name", vlan.Name)
	if _, err := client.Exchange(vlan.Name); err != nil {
		return fmt.Errorf("failed to exchange DHCP on %s: %w", vlan.Name, err)
	}
	vlanTestLogger.Info("succeeded to exchange DHCP", "name", vlan.Name)

	return nil
}
