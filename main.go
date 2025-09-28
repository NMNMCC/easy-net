package main

import (
	"github.com/alecthomas/kong"
	"nmnm.cc/easy-net/cli"
	"nmnm.cc/easy-net/internal/auth"
	"nmnm.cc/easy-net/internal/vlan"
)

var CLI struct {
	Auth cli.AuthCLI `cmd:"" help:"Authentication commands."`
	Vlan cli.VlanCLI `cmd:"" help:"VLAN commands."`
}

func main() {
	k := kong.Parse(&CLI)

	switch k.Command() {
	case "auth login":
		if CLI.Auth.Login.Base == *new(string) {
			base, err := auth.FindPortal(CLI.Auth.Host)
			if err != nil {
				k.FatalIfErrorf(err)
			}
			CLI.Auth.Login.Base = base
		}
		if err := auth.Login(&auth.LoginConfig{
			Base:     CLI.Auth.Login.Base,
			Link:     CLI.Auth.Link,
			UserID:   CLI.Auth.Login.UserID,
			Password: CLI.Auth.Login.Password,
		}); err != nil {
			k.FatalIfErrorf(err)
		}
	case "auth logout":
		if CLI.Auth.Login.Base == *new(string) {
			base, err := auth.FindPortal(CLI.Auth.Host)
			if err != nil {
				k.FatalIfErrorf(err)
			}
			CLI.Auth.Login.Base = base
		}
		if err := auth.Logout(&auth.LogoutConfig{
			Base:   CLI.Auth.Logout.Base,
			Link:   CLI.Auth.Link,
			UserID: CLI.Auth.Logout.UserID,
		}); err != nil {
			k.FatalIfErrorf(err)
		}
	case "auth attack":
		if err := auth.Attack(&auth.AttackConfig{
			Host:     CLI.Auth.Host,
			Link:     CLI.Auth.Link,
			Password: CLI.Auth.Attack.Password,
		}); err != nil {
			k.FatalIfErrorf(err)
		}
	case "vlan attack":
		if err := vlan.Attack(&vlan.AttackConfig{
			Start: CLI.Vlan.Attack.Start,
			Link:  CLI.Vlan.Link,
		}); err != nil {
			k.FatalIfErrorf(err)
		}
	case "vlan test":
		if err := vlan.Test(&vlan.TestConfig{
			Link: CLI.Vlan.Link,
			ID:   CLI.Vlan.Test.ID,
		}); err != nil {
			k.FatalIfErrorf(err)
		}
	default:
		k.PrintUsage(true)
		panic(k.Command())
	}
}
