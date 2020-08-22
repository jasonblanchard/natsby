package natsby

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusCollector observer for sending metrics to Prometheus
type PrometheusCollector struct {
	port                    string
	messagesReceivedCounter *prometheus.CounterVec
	repliesSentCounter      *prometheus.CounterVec
	latencyHistogram        *prometheus.HistogramVec
	errorsCounter           *prometheus.CounterVec
	isTesting               bool
}

// NewPrometheusCollector initialize new Prometheus observer
func NewPrometheusCollector(port string) *PrometheusCollector {
	p := &PrometheusCollector{}
	p.setupCollectors()
	p.port = port
	return p
}

// Collect starts prometheus server on p.port in a new goroutine
func (p *PrometheusCollector) Collect() error {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if p.isTesting == false {
			http.ListenAndServe(":"+p.port, nil)
		}
	}()
	return nil
}

// CollectSubjectReceived increment counter when message is received
func (p *PrometheusCollector) CollectSubjectReceived(subject string) {
	p.messagesReceivedCounter.WithLabelValues(subject).Inc()
}

// CollectLatency set histogram value for latency
func (p *PrometheusCollector) CollectLatency(subject string, latency time.Duration) {
	p.latencyHistogram.WithLabelValues(subject).Observe(latency.Seconds())
}

// CollectReply increment counter for replies
func (p *PrometheusCollector) CollectReply(subject string) {
	p.repliesSentCounter.WithLabelValues(subject).Inc()
}

// CollectError increment counter for errors
func (p *PrometheusCollector) CollectError(subject string) {
	p.errorsCounter.WithLabelValues(subject).Inc()
}

func (p *PrometheusCollector) setupCollectors() error {
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
