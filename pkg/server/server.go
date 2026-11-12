package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go1f/pkg/api"

	"github.com/go-chi/chi/v5"
)

const webDir = "./web"
const dateFormat = "20060102"

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var now time.Time
	var err error

	nowStr := queryParams.Get("now")
	if nowStr == "" {
		t := time.Now()
		now = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	date := queryParams.Get("date")
	repeat := queryParams.Get("repeat")

	res, err := api.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(res))
}

func Run(port string) error {
	r := chi.NewRouter()

	r.Get("/api/nextdate", nextDateHandler)
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	serverErr := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Сервер слушает порт: %s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("listen and serve: %w", err)
	case sig := <-shutdown:
		log.Printf("Got %s signal, shutting down", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
	}
	return nil
}
