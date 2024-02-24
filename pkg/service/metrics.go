package service

import (
	"net/http"

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

	err := http.ListenAndServe(":8081", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	if err != nil {
		logger.Fatal(err)
	}
}
