package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/willoma/keepakonf/internal"
)

const port = 35653

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
	srv, err := internal.Run(port)
	if err != nil {
		slog.Error("Could not start Keepakonf", "error", err)
		os.Exit(1)
	}

	<-done

	if err := srv.Close(); err != nil {
		slog.Error("Could not stop Keepakonf", "error", err)
		os.Exit(1)
	}
}
