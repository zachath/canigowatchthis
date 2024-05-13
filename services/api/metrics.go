package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "canigowatchthis_requests_total",
		Help: "Total number of requests by team",
	}, []string{"team"})

	errorsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "canigowatchthis_errors_total",
		Help: "Total number of errors by status code",
	}, []string{"statusCode"})
)
