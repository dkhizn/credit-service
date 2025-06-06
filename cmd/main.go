// Package main is the entry point for the credit-service application.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	httpAdapter "credit-service/internal/adapters/primary/http-adapter"
	creditService "credit-service/internal/application/credit-service"
	"credit-service/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Printf("failed to load config: %v", err)
		return
	}

	creditService := creditService.NewCreditService()
	httpAdapter, err := httpAdapter.New(log.Default(), cfg, creditService)
	if err != nil {
		log.Printf("failed to create HTTP adapter: %v", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := httpAdapter.Start(ctx); err != nil {
		log.Printf("HTTP server error: %v", err)
		return
	}
}
