package main

import (
	"regexp"
	"strconv"
	"time"

	"github.com/juju/errors"

	"github.com/amyangfei/go-logster/inter"
)

// LogReg is a simple http log regex
var LogReg = regexp.MustCompile(`.*HTTP/1.\d" (?P<http_status_code>\d{3}) .*`)

// SampleParser holds http status result
type SampleParser struct {
	http1xx     int
	http2xx     int
	http3xx     int
	http4xx     int
	http5xx     int
	httpUnknown int
}

// Init inits a *SampleParser Parser
func (parser *SampleParser) Init(options string) error {
	parser.http1xx = 0
	parser.http2xx = 0
	parser.http3xx = 0
	parser.http4xx = 0
	parser.http5xx = 0
	parser.httpUnknown = 0

	return nil
}

// ParseLine parses one line http log and caches parsed result
func (parser *SampleParser) ParseLine(line string) error {
	match := LogReg.FindStringSubmatch(line)
	if len(match) != LogReg.NumSubexp()+1 {
		return errors.Errorf("regex failed to match: %s", line)
	}
	status, err := strconv.Atoi(match[1])
	if err != nil {
		return errors.Trace(err)
	}
	switch {
	case status < 200:
		parser.http1xx++
	case status < 300:
		parser.http2xx++
	case status < 400:
		parser.http3xx++
	case status < 500:
		parser.http4xx++
	case status < 600:
		parser.http5xx++
	default:
		parser.httpUnknown++
	}
	return nil
}

// GetState gets http status metrics from cached parsed result
func (parser *SampleParser) GetState(duration float64) ([]*inter.Metric, error) {
	units := "Responses per sec"
	now := time.Now().Unix()
	return []*inter.Metric{
		{Name: "http_1xx", Value: float64(parser.http1xx) / duration, Units: units, Timestamp: now},
		{Name: "http_2xx", Value: float64(parser.http2xx) / duration, Units: units, Timestamp: now},
		{Name: "http_3xx", Value: float64(parser.http3xx) / duration, Units: units, Timestamp: now},
		{Name: "http_4xx", Value: float64(parser.http4xx) / duration, Units: units, Timestamp: now},
		{Name: "http_5xx", Value: float64(parser.http5xx) / duration, Units: units, Timestamp: now},
		{Name: "http_unknown", Value: float64(parser.httpUnknown) / duration, Units: units, Timestamp: now},
	}, nil
}

func main() {}

// Parser declares a SampleParser object
var Parser SampleParser
