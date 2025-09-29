package tool

import (
	"fmt"
	"io"
	"time"

	"github.com/dustin/go-humanize"
	"nmnm.cc/easy-net/internal/log"
	"nmnm.cc/easy-net/internal/util"
)

type SpeedtestConfig struct {
	URL     string
	Link    string
	Timeout time.Duration
}

var speedtestLogger = log.New("tool/speedtest")

func Speedtest(cfg *SpeedtestConfig) error {
	client := util.NewHTTPClient(cfg.Link)

	speedtestLogger.Info("performing speed test", "url", cfg.URL)
	res, err := client.Get(cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to perform speed test: %w", err)
	}
	defer res.Body.Close()

	go func() {
		time.Sleep(cfg.Timeout)
		res.Body.Close()
	}()

	start := time.Now()
	written, _ := io.Copy(io.Discard, res.Body)
	duration := time.Since(start)

	speed := humanize.Bytes(uint64(float64(written)/duration.Seconds())) + "/s"

	speedtestLogger.Info("speed test completed", "url", cfg.URL, "duration", duration, "bytes", written, "speed", speed)

	return nil
}
