//go:build linux || android
// +build linux android

package vlan

import (
	"log/slog"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/vishvananda/netlink"
)

type TestConfig struct {
	Link string
	ID   uint16
}

func Test(cfg *TestConfig) error {
	logger := slog.With("component", "vlan", "id", cfg.ID)

	master, err := netlink.LinkByName(cfg.Link)
	if err != nil {
		return err
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = NewVlanName(cfg.Link, cfg.ID)
	attrs.ParentIndex = master.Attrs().Index

	vlan := &netlink.Vlan{
		LinkAttrs: attrs,
		VlanId:    int(cfg.ID),
	}

	logger.Info("adding VLAN interface", "link", cfg.Link)
	if err := netlink.LinkAdd(vlan); err != nil {
		logger.Error("failed to add VLAN interface", "name", vlan.Name, "error", err)
		return err
	}
	defer func() {
		logger.Info("removing VLAN interface", "name", vlan.Name)
		if err := netlink.LinkDel(vlan); err != nil {
			logger.Error("failed to remove VLAN interface", "name", vlan.Name, "error", err)
			return
		}
		logger.Info("removed VLAN interface", "name", vlan.Name)
	}()
	logger.Info("added VLAN interface", "name", vlan.Name)

	// logger.Info("setting master and bringing up", "name", vlan.Name)
	// if err := netlink.LinkSetMaster(vlan, master); err != nil {
	// 	logger.Error("failed to set master for VLAN interface", "name", vlan.Name, "error", err)
	// 	return err
	// }
	// logger.Info("set master for VLAN interface", "name", vlan.Name)

	logger.Info("bringing up VLAN interface", "name", vlan.Name)
	if err := netlink.LinkSetUp(vlan); err != nil {
		logger.Error("failed to bring up VLAN interface", "name", vlan.Name, "error", err)
		return err
	}
	logger.Info("brought up VLAN interface", "name", vlan.Name)

	client := &client4.Client{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Info("exchanging DHCP", "name", vlan.Name)
	if _, err := client.Exchange(vlan.Name); err != nil {
		logger.Error("failed to exchange DHCP", "name", vlan.Name, "error", err)
		return err
	}
	logger.Info("succeeded to exchange DHCP", "name", vlan.Name)

	return nil
}
