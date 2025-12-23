package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"log"

	"github.com/fox-md/fox-dummy-keda-scaler/scaler"
)

func StartHTTP(ctx context.Context, addr string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/get", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(scaler.DummyCounter)))
	})

	mux.HandleFunc("GET /up/{increment}", func(w http.ResponseWriter, r *http.Request) {

		inc, err := strconv.Atoi(r.PathValue("increment"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Value must be an integer"))
		}

		scaler.DummyCounter += inc

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /down/{decrement}", func(w http.ResponseWriter, r *http.Request) {

		dec, err := strconv.Atoi(r.PathValue("decrement"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Value must be an integer"))
		}

		if scaler.DummyCounter-dec < 1 {
			scaler.DummyCounter = 1
		} else {
			scaler.DummyCounter -= dec
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Print("stopping http server")
		ctxShutdown, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()
		srv.Shutdown(ctxShutdown)
	}()

	log.Printf("http listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http serve: %v", err)
	}
}
