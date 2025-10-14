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

var Mb = uint64(1_000_000 / 8)

type SpeedtestResult struct {
	Speed uint64
}

var speedtestLogger = log.New("tool/speedtest")

func Speedtest(cfg *SpeedtestConfig) (*SpeedtestResult, error) {
	client := util.NewHTTPClient(cfg.Link)

	speedtestLogger.Info("performing speed test", "url", cfg.URL)
	res, err := client.Get(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to perform speed test: %w", err)
	}
	defer res.Body.Close()

	go func() {
		time.Sleep(cfg.Timeout)
		res.Body.Close()
	}()

	start := time.Now()
	written, _ := io.Copy(io.Discard, res.Body)
	duration := time.Since(start)

	speed := uint64(float64(written) / duration.Seconds())
	speed_human := humanize.Bytes(speed) + "/s"

	speedtestLogger.Info("speed test completed", "url", cfg.URL, "duration", duration, "bytes", written, "speed", speed_human)

	return &SpeedtestResult{
		Speed: speed,
	}, nil
}
