package line

import (
	"crypto/rand"
	"fmt"
	"net"
	"strings"
	"time"

	"nmnm.cc/easy-net/internal/log"
)

var morseCodeMap = map[rune]string{
	'A': "._",
	'B': "_...",
	'C': "_._.",
	'D': "_..",
	'E': ".",
	'F': ".._.",
	'G': "__.",
	'H': "....",
	'I': "..",
	'J': ".___",
	'K': "_._",
	'L': "._..",
	'M': "__",
	'N': "_.",
	'O': "___",
	'P': ".__.",
	'Q': "__._",
	'R': "._.",
	'S': "...",
	'T': "_",
	'U': ".._",
	'V': "..._",
	'W': ".__",
	'X': "_.._",
	'Y': "_.__",
	'Z': "__..",

	'1': ".____",
	'2': "..___",
	'3': "...__",
	'4': "...._",
	'5': ".....",
	'6': "_....",
	'7': "__...",
	'8': "___..",
	'9': "____.",
	'0': "_____",

	'.':  "._._._",
	',':  "__..__",
	'?':  "..__..",
	'\'': ".____.",
	'!':  "_._.__",
	'/':  "_.._.",
	'(':  "_.__.",
	')':  "_.__._",
	'&':  "._...",
	':':  "___...",
	';':  "_._._.",
	'=':  "_..._",
	'+':  "._._.",
	'-':  "_...._",
	'_':  "..__._",
	'"':  "._.._.",
	'$':  "..._.._",
	'@':  ".__._.",
}

type SendMorseMessageConfig struct {
	Interval time.Duration
	Message  string
	Times    uint64
}

var sendMorseMessageLogger = log.New("line/morse")

func SendMorseMessage(cfg *SendMorseMessageConfig) error {
	sendMorseMessageLogger.Info("sending morse message", "message", cfg.Message)

	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:0")
	if err != nil {
		sendMorseMessageLogger.Error("failed to resolve udp addr", "err", err)
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		sendMorseMessageLogger.Error("failed to dial udp", "err", err)
		return err
	}
	defer conn.Close()

	for i := range cfg.Times {
		sendMorseMessageLogger.Info("sending message", "round", i+1, "total", cfg.Times)

		for _, c := range strings.ToUpper(cfg.Message) {
			if c == ' ' {
				fmt.Println()
				time.Sleep(cfg.Interval * 7)
				continue
			}

			morse, ok := morseCodeMap[c]
			if !ok {
				sendMorseMessageLogger.Warn("skipping unsupported character", "char", c)
				continue
			}

			fmt.Printf("%c ", c)

			for _, s := range morse {
				fmt.Printf("%c", s)

				size := (1 << (10 * 2))

				switch s {
				case '.':
					for range size {
						buf := make([]byte, 1)
						if _, err := rand.Read(buf); err != nil {
							sendMorseMessageLogger.Error("failed to read random bytes", "err", err)
							return err
						}

						conn.WriteTo(buf, addr)
					}
				case '_':
					for range size * 10 {
						buf := make([]byte, 1)
						if _, err := rand.Read(buf); err != nil {
							sendMorseMessageLogger.Error("failed to read random bytes", "err", err)
							return err
						}

						conn.WriteTo(buf, addr)
					}
				}

				time.Sleep(cfg.Interval * 3)
			}

			fmt.Println()
		}

		sendMorseMessageLogger.Info("finished sending message", "round", i+1, "total", cfg.Times)
	}

	sendMorseMessageLogger.Info("morse message sent successfully")
	return nil
}
