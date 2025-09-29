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
		Timeout  time.Duration `help:"Timeout for connection verification." default:"15s" short:"t"`
		// TargetSpeed string        `help:"Target speed for connection verification." default:"15Mbps" short:"s"`
	} `cmd:"" help:"Attack the network."`
}
