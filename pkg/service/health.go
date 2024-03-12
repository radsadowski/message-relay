package service

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func healthCheckHandler(healthCh chan<- struct{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		healthCh <- struct{}{}
		w.WriteHeader(http.StatusNoContent)
	}
}

func startHealthCheckServer(ctx context.Context, healthCh chan<- struct{}, logger *log.Logger) {
	logger.Info("starting the health check server on :8080 ...")
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           http.HandlerFunc(healthCheckHandler(healthCh)),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       10 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("health check server shutdown failed:%+v", err)
	} else {
		log.Print("health check server exited properly")
	}
}
