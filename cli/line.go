package cli

import "time"

type LineCLI struct {
	Interval time.Duration `help:"Interval between Morse code signals." default:"10ms" short:"i"`
	Message  string        `arg:"" help:"Message to send." required:""`
	Times    uint64        `help:"Number of times to send the message." default:"999999999" short:"t"`
}
