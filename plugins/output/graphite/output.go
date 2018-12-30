package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/amyangfei/go-logster/logster"
	"github.com/buger/jsonparser"
	"github.com/juju/errors"
	"github.com/rs/zerolog"
)

// DefaultSeparator is separator used in metric operation
const DefaultSeparator = "."

// GraphiteOutput sends metric to graphite
type GraphiteOutput struct {
	logster.MetricOp
	Host     string
	Prototol string
	DryRun   bool
	Logger   zerolog.Logger
}

func parserKey(options, key, defaultVal string) (string, error) {
	val, dataType, _, err := jsonparser.Get([]byte(options), key)
	if err != nil {
		if defaultVal != "" && dataType == jsonparser.NotExist {
			return defaultVal, nil
		}
		return "", errors.Trace(err)
	}
	return string(val), nil
}

// Init inits the *GraphiteOutput type Output
func (output *GraphiteOutput) Init(
	prefix, suffix, options string, dryRun bool, logger zerolog.Logger) error {
	host, err := parserKey(options, "host", "")
	if err != nil {
		return errors.Trace(err)
	}
	protocol, err := parserKey(options, "protocol", "udp")
	if err != nil {
		return errors.Trace(err)
	}

	separator, err := parserKey(options, "separator", DefaultSeparator)
	if err != nil {
		return errors.Trace(err)
	}

	output.Host = host
	output.Prototol = protocol
	output.DryRun = dryRun
	output.Logger = logger
	output.MetricOp.Separator = separator
	output.MetricOp.Prefix = prefix
	output.MetricOp.Suffix = suffix
	return nil
}

// Submit sends metrics to graphite
func (output *GraphiteOutput) Submit(metrics []*logster.Metric) error {
	var conn net.Conn
	var err error
	if !output.DryRun {
		conn, err = net.Dial(output.Prototol, output.Host)
		if err != nil {
			return errors.Trace(err)
		}
		defer conn.Close()
	}
	for _, metric := range metrics {
		metricName := output.MetricOp.GetMetricName(metric)
		if strings.Contains(metricName, " ") {
			return errors.Errorf("Invalid metric name: \"%s\", spaces not allowed", metricName)
		}
		mstr := fmt.Sprintf("%s %v %d", metricName, metric.Value, metric.Timestamp)
		output.Logger.Debug().Msgf("submitting graphite metric: %s", mstr)
		if output.DryRun {
			fmt.Printf("%s %s\n", output.Host, mstr)
		} else {
			conn.Write([]byte(mstr))
		}
	}
	return nil
}

func main() {}

// Output declares an GraphiteOutput object
var Output GraphiteOutput
