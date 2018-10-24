package logster

type Parser interface {
	Init(options string) error
	ParseLine(line string) error
	GetState(duration float64) ([]*Metric, error)
}

type Metric struct {
	Name       string
	Value      float64
	Units      string
	Timestamp  int64
	MetricType string
}

type Output interface {
	Init(prefix, suffix, options string) error
	Submit(metrics []*Metric) error
}

type MetricOp struct {
	Prefix    string
	Suffix    string
	Separator string
}

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
