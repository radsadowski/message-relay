package service

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var relayedMessagesCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "relayed_messages_total",
	Help: "The total number of relayed messages",
})

func startMetricsServer(logger *log.Logger) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(relayedMessagesCounter)
	logger.Info("starting the metrics server on :8081 ...")
	server := &http.Server{
		Addr:              ":8081",
		Handler:           promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
