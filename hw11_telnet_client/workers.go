package main

import (
	"context"
	"fmt"
	"os"
)

func serve(ctx context.Context, tc TelnetClient) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := tc.Send(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
		}
	}
}

func listen(ctx context.Context, tc TelnetClient, cancel context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := tc.Receive()
			if err != nil {
				fmt.Fprintln(os.Stderr, "...Connection closed by peer")
				cancel()
				return
			}
		}
	}
}
