package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	pvzCounter       prometheus.Counter
	receptionCounter prometheus.Counter
	productsCounter  prometheus.Counter
	userCounter      prometheus.Counter
	totalCounter     prometheus.Counter
	httpDuration     prometheus.Histogram
}

func NewMetrics() *Metrics {
	return &Metrics{
		pvzCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "avito.pvz.pvzs_total",
			Help: "The total number of created pvzs",
		}),
		receptionCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "avito.pvz.receptions_total",
			Help: "The total number of created receptions",
		}),
		productsCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "avito.pvz.products_total",
			Help: "The total number of created products",
		}),
		userCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "avito.pvz.user_total",
			Help: "The total number of created users",
		}),
		totalCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "avito.pvz.requests_total",
			Help: "The total number of requests",
		}),
		httpDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "avito.pvz.requests_duration",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		}),
	}
}

func (m Metrics) IncPVZs() {
	m.pvzCounter.Inc()
}

func (m Metrics) IncReceptions() {
	m.receptionCounter.Inc()
}

func (m Metrics) IncProducts() {
	m.productsCounter.Inc()
}

func (m Metrics) IncUsers() {
	m.userCounter.Inc()
}

func (m Metrics) IncRequests(duration time.Duration) {
	m.httpDuration.Observe(duration.Seconds())
	m.totalCounter.Inc()
}
