package inter

import (
	"github.com/rs/zerolog"
)

// Metric holds log analysis result set
type Metric struct {
	Name      string
	Value     interface{}
	Units     string
	Timestamp int64
}

// MetricOp defines some common operation to metric
type MetricOp struct {
	Prefix    string
	Suffix    string
	Separator string
}

// Parser defines an interface to a log parser
type Parser interface {
	Init(options string) error
	ParseLine(line string) error
	GetState(duration float64) ([]*Metric, error)
}

// Output defines an interface of metric sending target
type Output interface {
	Init(prefix, suffix, options string, dryRun bool, logger zerolog.Logger) error
	Submit(metrics []*Metric) error
}

// GetMetricName returns the operated name of a metric, the operation includes
// adding prefix and adding suffix now
func (op *MetricOp) GetMetricName(metric *Metric) string {
	metricName := metric.Name
	if op.Prefix != "" {
		metricName = op.Prefix + op.Separator + metricName
	}
	if op.Suffix != "" {
		metricName = metricName + op.Separator + op.Suffix
	}
	return metricName
}
