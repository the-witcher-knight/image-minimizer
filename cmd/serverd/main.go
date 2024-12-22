package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	v1 "github.com/the-witcher-knight/image-minimize-go/internal/handler/rest/v1"
	"github.com/the-witcher-knight/image-minimize-go/internal/service/imaging"
)

const (
	port          = "8080"
	serverTimeout = 5 * time.Second
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Printf("server exit abnormally %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Setup dependencies
	imgSvc, cleanup := imaging.New()
	defer cleanup()

	handler := v1.New(imgSvc)

	// Listen INTERRUPT event for graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// Setup http server
	srv := http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", port),
		Handler: routes(handler),
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Printf("web server started; listening at %s\n", srv.Addr)
		defer fmt.Println("web server shutdown")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), serverTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}

func routes(v1Handler v1.Handler) http.Handler {
	r := chi.NewRouter()
	r.Get("/public/healthz", livenessHandlerFunc)
	r.Route("/api/v1", func(v1Router chi.Router) {
		v1Router.Post("/image/resize", v1Handler.ReduceImageSize())
	})

	return r
}

// livenessHandlerFunc is a simple handler func to be used fo health check operations
func livenessHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = fmt.Fprintln(w, "ok")
}
