package service

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func healthCheckHandler(healthCh chan<- struct{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		healthCh <- struct{}{}
		w.WriteHeader(http.StatusOK)
	}
}

func startHealthCheckServer(healthCh chan<- struct{}, logger *log.Logger) {
	logger.Info("starting the health check server on :8080 ...")
	err := http.ListenAndServe(":8080", http.HandlerFunc(healthCheckHandler(healthCh)))
	if err != nil {
		logger.Fatal(err)
	}
}
