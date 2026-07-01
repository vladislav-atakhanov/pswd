package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func trapSignals() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		fmt.Fprintln(os.Stderr, "\ninterrupted")
		cancel()
	}()
	return ctx
}
