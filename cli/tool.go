package cli

import "time"

type ToolCLI struct {
	Morse struct {
		Interval time.Duration `help:"Interval between Morse code signals." default:"10ms" short:"i"`
		Message  string        `arg:"" help:"Message to send." required:""`
		Times    uint64        `help:"Number of times to send the message." default:"999999999" short:"t"`
	} `cmd:"" help:"Send a Morse code message over the network."`

	Speedtest struct {
		Link    string        `help:"Network link to use for speed test." short:"l"`
		URL     string        `help:"URL to test against." default:"https://yun.mcloud.139.com/hongseyunpan/2.43G.zip" short:"u"`
		Timeout time.Duration `help:"Timeout for the speed test." default:"5s" short:"t"`
	} `cmd:"" help:"Perform a speed test."`
}
