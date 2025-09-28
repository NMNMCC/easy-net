package main

import (
	"github.com/alecthomas/kong"
	"nmnm.cc/easy-net/cli"
	"nmnm.cc/easy-net/internal/auth"
	"nmnm.cc/easy-net/internal/log"
	"nmnm.cc/easy-net/internal/vlan"
)

var CLI struct {
	Auth cli.AuthCLI `cmd:"" help:"Authentication commands."`
	Vlan cli.VlanCLI `cmd:"" help:"VLAN commands."`
}

var logger = log.New("main")

func main() {
	k := kong.Parse(&CLI)

	var err error

	switch k.Command() {
	case "auth login":
		if CLI.Auth.Login.Base == *new(string) {
			var base string
			base, err = auth.FindPortal(CLI.Auth.Host)
			if err == nil {
				CLI.Auth.Login.Base = base
			}
		}
		if err == nil {
			err = auth.Login(&auth.LoginConfig{
				Base:     CLI.Auth.Login.Base,
				Link:     CLI.Auth.Link,
				UserID:   CLI.Auth.Login.UserID,
				Password: CLI.Auth.Login.Password,
			})
		}
	case "auth logout":
		if CLI.Auth.Logout.Base == *new(string) {
			var base string
			base, err = auth.FindPortal(CLI.Auth.Host)
			if err == nil {
				CLI.Auth.Logout.Base = base
			}
		}
		if err == nil {
			err = auth.Logout(&auth.LogoutConfig{
				Base:   CLI.Auth.Logout.Base,
				Link:   CLI.Auth.Link,
				UserID: CLI.Auth.Logout.UserID,
			})
		}
	case "auth attack":
		err = auth.Attack(&auth.AttackConfig{
			Host:     CLI.Auth.Host,
			Link:     CLI.Auth.Link,
			Password: CLI.Auth.Attack.Password,
		})
	case "vlan attack":
		err = vlan.Attack(&vlan.AttackConfig{
			Start: CLI.Vlan.Attack.Start,
			Link:  CLI.Vlan.Link,
		})
	case "vlan test":
		err = vlan.Test(&vlan.TestConfig{
			Link: CLI.Vlan.Link,
			ID:   CLI.Vlan.Test.ID,
		})
	default:
		k.PrintUsage(true)
		panic(k.Command())
	}

	if err != nil {
		logger.Error("command failed", "command", k.Command(), "error", err)
		k.Exit(1)
	}
}
