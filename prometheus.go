package natsby

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds metrics counters
type Metrics struct {
	subscriptionCounter     *prometheus.CounterVec
	messagesReceivedCounter *prometheus.CounterVec
	repliesSentCounter      *prometheus.CounterVec
	latencyHistogram        *prometheus.HistogramVec
	errorsCounter           *prometheus.CounterVec
}

// WithPrometheusInput configuration for prometheus middleware
type WithPrometheusInput struct {
	Port string
}

// WithPrometheus Records some metrics in prometheus format
func WithPrometheus(input *WithPrometheusInput) HandlerFunc {
	m := &Metrics{}

	m.messagesReceivedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "natsby_messages_received_total_by_subject",
			Help: "How many NATS messages received, partitioned by subject.",
		},
		[]string{"subject"},
	)

	m.repliesSentCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "natsby_replies_sent_total_by_subject",
			Help: "How many NATS message replies sent, partitioned by subject.",
		},
		[]string{"subject"},
	)

	m.latencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "natsby_latency_histogram",
		Help: "Histogram of latencies for message handling, partitioned by subject",
		// TODO: Make this configurable
		// Buckets: prometheus.LinearBuckets(20, 5, 5),
	},
		[]string{"subject"},
	)

	m.errorsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "natsby_errors_by_subject",
			Help: "Errors encountered, partitioned by subject.",
		},
		[]string{"subject"},
	)

	err := prometheus.Register(m.messagesReceivedCounter)
	err = prometheus.Register(m.repliesSentCounter)
	err = prometheus.Register(m.latencyHistogram)
	err = prometheus.Register(m.errorsCounter)

	if err != nil {
		// TODO: Do something logical, here
		panic(err)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":"+input.Port, nil)
	}()

	return func(c *Context) {
		m.messagesReceivedCounter.WithLabelValues(c.Msg.Subject).Inc()
		start := time.Now()

		c.Next()

		_, ok := c.GetByteReplyPayload()

		if ok == true {
			m.repliesSentCounter.WithLabelValues(c.Msg.Subject).Inc()
		}

		end := time.Now()
		latency := end.Sub(start) // TODO: Duplicating latency calculation right now in logger
		m.latencyHistogram.WithLabelValues(c.Msg.Subject).Observe(latency.Seconds())

		if c.Err != nil {
			m.errorsCounter.WithLabelValues(c.Msg.Subject).Inc()
		}
	}
}
