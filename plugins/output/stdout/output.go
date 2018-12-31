package main

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/juju/errors"
	"github.com/rs/zerolog"

	"github.com/amyangfei/go-logster/inter"
)

// DefaultSeparator is separator used in metric operation
const DefaultSeparator = "."

// StdoutOutput sends metrics to stdout
type StdoutOutput struct {
	inter.MetricOp
}

// Init inits the *StdoutOutput type Output
func (output *StdoutOutput) Init(
	prefix, suffix, options string, dryRun bool, logger zerolog.Logger) error {
	val, dataType, _, err := jsonparser.Get([]byte(options), "separator")
	if err != nil {
		if dataType == jsonparser.NotExist {
			output.MetricOp.Separator = DefaultSeparator
		} else {
			return errors.Trace(err)
		}
	} else {
		output.MetricOp.Separator = string(val)
	}
	output.MetricOp.Prefix = prefix
	output.MetricOp.Suffix = suffix
	return nil
}

// Submit send metrics to stdout
func (output *StdoutOutput) Submit(metrics []*inter.Metric) error {
	for _, metric := range metrics {
		metricName := output.MetricOp.GetMetricName(metric)
		fmt.Printf("%d %s %v\n", metric.Timestamp, metricName, metric.Value)
	}
	return nil
}

func main() {}

// Output declares a StdoutOutput object
var Output StdoutOutput
