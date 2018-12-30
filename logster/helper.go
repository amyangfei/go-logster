package logster

import (
	"plugin"

	"github.com/juju/errors"
	"github.com/rs/zerolog"
)

// Parser defines an interface to a log parser
type Parser interface {
	Init(options string) error
	ParseLine(line string) error
	GetState(duration float64) ([]*Metric, error)
}

// Metric holds log analysis result set
type Metric struct {
	Name      string
	Value     interface{}
	Units     string
	Timestamp int64
}

// Output defines an interface of metric sending target
type Output interface {
	Init(prefix, suffix, options string, dryRun bool, logger zerolog.Logger) error
	Submit(metrics []*Metric) error
}

// MetricOp defines some common operation to metric
type MetricOp struct {
	Prefix    string
	Suffix    string
	Separator string
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

// LoadParserPlugin loads a parser plugin from given plugin path,
// the plugin name is specificed to Parser
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

// LoadOutputPlugin loads a Output plugin from given plugin path,
// the plugin name is specificed to Output
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
