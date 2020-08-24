package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const usage = "usage: go-telnet --timeout=10s host port"

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Second*10, usage)
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatalf("not enough arguments %s\n", usage)
	}

	tc := NewTelnetClient(net.JoinHostPort(flag.Arg(0), flag.Arg(1)), timeout, os.Stdin, os.Stdout)

	err := tc.Connect()
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	defer tc.Close()
	fmt.Fprintln(os.Stderr, "...Connected to", net.JoinHostPort(flag.Arg(0), flag.Arg(1)))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()

	go listen(ctx, tc, cancel)
	serve(ctx, tc)
}
