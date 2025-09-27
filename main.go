package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"nmnm.cc/easy-net/internal"
)

const (
	DefaultPassword = "112233"
)

func main() {
	if internal.TestConnection() {
		return
	}

	base, err := internal.FindPortal("3.3.3.3")
	if err != nil {
		fmt.Printf("failed to find portal: %s\n", err)
		return
	}
	slog.Info("portal base", "base", base)

	for {
		tried := make(map[string]struct{})

		var userid string
		for {
			userid = internal.RandomUserid()

			if _, ok := tried[userid]; ok {
				continue
			} else {
				tried[userid] = struct{}{}
			}

			err := internal.Login(base, userid, DefaultPassword)
			if err != nil {
				fmt.Printf("login failed: %s %s\n", userid, err)
				continue
			}

			fmt.Printf("login succeeded: %s\n", userid)
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		go func() {
			for {
				time.Sleep(1 * time.Second)

				if internal.TestConnection() {
					cancel()
				}
			}
		}()

		<-ctx.Done()

		switch ctx.Err() {
		case context.DeadlineExceeded:
			internal.Logout(base, userid)
			fmt.Printf("login timed out: %s\n", userid)
		case context.Canceled:
			return
		}
	}
}
