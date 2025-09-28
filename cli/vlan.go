package cli

type VlanCLI struct {
	Link string `help:"Network link to use for VLAN attack." short:"l"`
	Test struct {
		ID uint16 `help:"VLAN ID to try." required:"" short:"i"`
	} `cmd:"try" help:"Test if a VLAN ID is valid."`
	Attack struct {
		Start uint16 `help:"Starting VLAN ID to try." default:"1586" short:"s"`
	} `cmd:"" help:"Perform VLAN attack."`
}
