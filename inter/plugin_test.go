package inter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricName(t *testing.T) {
	mo := &MetricOp{
		Prefix:    "server.101",
		Suffix:    "count",
		Separator: "_",
	}
	metric := &Metric{
		Name:      "login",
		Value:     1000,
		Timestamp: time.Now().Unix(),
	}
	name := mo.GetMetricName(metric)
	assert.Equal(t, "server.101_login_count", name)
}
