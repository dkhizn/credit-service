// Package httpadapter реализует HTTP-адаптер кредитного сервиса.
package httpadapter

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"credit-service/internal/adapters/primary/http-adapter/controller"
	"credit-service/internal/adapters/primary/http-adapter/middleware"
	"credit-service/internal/adapters/primary/http-adapter/router"
	creditService "credit-service/internal/application/credit-service"
	"credit-service/internal/config"
)

// HTTPAdapter предоставляет методы запуска HTTP-сервера.
type HTTPAdapter struct {
	server *http.Server
	logger *log.Logger
}

// New создаёт новый HTTP-адаптер с логгером, конфигом и сервисом.
func New(logger *log.Logger, cfg *config.Config, creditSvc *creditService.CredtiService) (*HTTPAdapter, error) {
	ctr := controller.New(creditSvc)
	r := router.NewRouter()
	r.RegisterRoutes(ctr)
	handlerWithMiddleware := middleware.LoggingMiddleware(r.Router())

	srv := &http.Server{
		Addr:         cfg.HTTP.Port,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &HTTPAdapter{
		server: srv,
		logger: logger,
	}, nil
}

// Start запускает HTTP-сервер и слушает системные сигналы завершения.
func (h *HTTPAdapter) Start(ctx context.Context) error {
	h.logger.Printf("HTTP server listening on %s\n", h.server.Addr)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		h.logger.Println("Shutting down HTTP server…")
		if err := h.server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		err := h.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
