package natsby

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusObserver observer for sending metrics to Prometheus
type PrometheusObserver struct {
	port                    string
	messagesReceivedCounter *prometheus.CounterVec
	repliesSentCounter      *prometheus.CounterVec
	latencyHistogram        *prometheus.HistogramVec
	errorsCounter           *prometheus.CounterVec
	isTesting               bool
}

// NewPrometheusObserver initialize new Prometheus observer
func NewPrometheusObserver(port string) *PrometheusObserver {
	p := &PrometheusObserver{}
	p.setupCollectors()
	p.port = port
	return p
}

// Collect starts prometheus server on p.port in a new goroutine
func (p *PrometheusObserver) Collect() error {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if p.isTesting == false {
			http.ListenAndServe(":"+p.port, nil)
		}
	}()
	return nil
}

// ObserveSubjectReceived increment counter when message is received
func (p *PrometheusObserver) ObserveSubjectReceived(subject string) {
	p.messagesReceivedCounter.WithLabelValues(subject).Inc()
}

// ObserveLatency set histogram value for latency
func (p *PrometheusObserver) ObserveLatency(subject string, latency time.Duration) {
	p.latencyHistogram.WithLabelValues(subject).Observe(latency.Seconds())
}

// ObserveReply increment counter for replies
func (p *PrometheusObserver) ObserveReply(subject string) {
	p.repliesSentCounter.WithLabelValues(subject).Inc()
}

// ObserveError increment counter for errors
func (p *PrometheusObserver) ObserveError(subject string) {
	p.errorsCounter.WithLabelValues(subject).Inc()
}

func (p *PrometheusObserver) setupCollectors() error {
	p.messagesReceivedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "natsby_messages_received_total_by_subject",
			Help: "How many NATS messages received, partitioned by subject.",
		},
		[]string{"subject"},
	)

	p.repliesSentCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "natsby_replies_sent_total_by_subject",
			Help: "How many NATS message replies sent, partitioned by subject.",
		},
		[]string{"subject"},
	)

	p.latencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "natsby_latency_histogram",
		Help: "Histogram of latencies for message handling, partitioned by subject",
		// TODO: Make this configurable
		// Buckets: prometheus.LinearBuckets(20, 5, 5),
	},
		[]string{"subject"},
	)

	p.errorsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "natsby_errors_by_subject",
			Help: "Errors encountered, partitioned by subject.",
		},
		[]string{"subject"},
	)

	err := prometheus.Register(p.messagesReceivedCounter)
	err = prometheus.Register(p.repliesSentCounter)
	err = prometheus.Register(p.latencyHistogram)
	err = prometheus.Register(p.errorsCounter)

	return err
}
