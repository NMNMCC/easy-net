package cli

import "time"

type AuthCLI struct {
	Host  string `help:"Host address." short:"H" default:"3.3.3.3"`
	Link  string `help:"Network link name." short:"l"`
	Login struct {
		Base     string `help:"Base URL of the authentication portal."`
		UserID   string `help:"User ID." required:"" short:"u"`
		Password string `help:"Password." required:"" short:"p"`
	} `cmd:"" help:"Login to the network."`
	Logout struct {
		Base   string `help:"Base URL of the authentication portal."`
		UserID string `help:"User ID." required:"" short:"u"`
	} `cmd:"" help:"Logout from the network."`
	Attack struct {
		Password string        `help:"Password." required:"" default:"112233" short:"p"`
		Wait     time.Duration `help:"Wait time after successful login to verify connection." default:"15s" short:"w"`
	} `cmd:"" help:"Attack the network."`
}
