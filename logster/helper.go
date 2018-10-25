package logster

import (
	"errors"
	"plugin"
)

type Parser interface {
	Init(options string) error
	ParseLine(line string) error
	GetState(duration float64) ([]*Metric, error)
}

type Metric struct {
	Name      string
	Value     interface{}
	Units     string
	Timestamp int64
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

func LoadParserPlugin(pluginPath string) (Parser, error) {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	symbol, err := plug.Lookup("Parser")
	if err != nil {
		return nil, err
	}
	parser, ok := symbol.(Parser)
	if !ok {
		return nil, errors.New("unexpected type from module symbol")
	}
	return parser, nil
}

func LoadOutputPlugin(pluginPath string) (Output, error) {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	symbol, err := plug.Lookup("Output")
	if err != nil {
		return nil, err
	}
	output, ok := symbol.(Output)
	if !ok {
		return nil, errors.New("unexpected type from module symbol")
	}
	return output, nil
}
