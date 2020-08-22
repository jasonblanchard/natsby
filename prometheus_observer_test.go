package natsby

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPrometheusObserver(t *testing.T) {
	p := NewPrometheusObserver("1234")
	p.isTesting = true

	var ptype *PrometheusObserver
	assert.IsType(t, ptype, p)

	p.ObserveSubjectReceived("test.subject")
	p.ObserveLatency("test.subject", 1*time.Second)
	p.ObserveReply("test.subject")
	p.ObserveError("test.subject")
	p.Collect()
}
