package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	IssuedOrdersCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "issued_orders_total",
		Help: "Total number of issued orders",
	})
)

func Init() {
	prometheus.MustRegister(IssuedOrdersCounter)
}
