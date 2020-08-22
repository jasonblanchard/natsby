package natsby

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPrometheusCollector(t *testing.T) {
	p := NewPrometheusCollector("1234")
	p.isTesting = true

	var ptype *PrometheusCollector
	assert.IsType(t, ptype, p)

	p.CollectSubjectReceived("test.subject")
	p.CollectLatency("test.subject", 1*time.Second)
	p.CollectReply("test.subject")
	p.CollectError("test.subject")
	p.Collect()
}
