package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/amyangfei/go-logster/logster"
)

var LogReg = regexp.MustCompile(`.*HTTP/1.\d" (?P<http_status_code>\d{3}) .*`)

type SampleParser struct {
	http1xx     int
	http2xx     int
	http3xx     int
	http4xx     int
	http5xx     int
	httpUnknown int
}

func (parser *SampleParser) Init(options string) error {
	parser.http1xx = 0
	parser.http2xx = 0
	parser.http3xx = 0
	parser.http4xx = 0
	parser.http5xx = 0
	parser.httpUnknown = 0

	return nil
}

func (parser *SampleParser) ParseLine(line string) error {
	match := LogReg.FindStringSubmatch(line)
	if len(match) != LogReg.NumSubexp()+1 {
		return fmt.Errorf("regex failed to match: %s", line)
	}
	status, err := strconv.Atoi(match[1])
	if err != nil {
		return err
	}
	switch {
	case status < 200:
		parser.http1xx += 1
	case status < 300:
		parser.http2xx += 1
	case status < 400:
		parser.http3xx += 1
	case status < 500:
		parser.http4xx += 1
	case status < 600:
		parser.http5xx += 1
	default:
		parser.httpUnknown += 1
	}
	return nil
}

func (parser *SampleParser) GetState(duration float64) ([]*logster.Metric, error) {
	units := "Responses per sec"
	now := time.Now().Unix()
	return []*logster.Metric{
		&logster.Metric{Name: "http_1xx", Value: float64(parser.http1xx) / duration, Units: units, Timestamp: now},
		&logster.Metric{Name: "http_2xx", Value: float64(parser.http2xx) / duration, Units: units, Timestamp: now},
		&logster.Metric{Name: "http_3xx", Value: float64(parser.http3xx) / duration, Units: units, Timestamp: now},
		&logster.Metric{Name: "http_4xx", Value: float64(parser.http4xx) / duration, Units: units, Timestamp: now},
		&logster.Metric{Name: "http_5xx", Value: float64(parser.http5xx) / duration, Units: units, Timestamp: now},
		&logster.Metric{Name: "http_unknown", Value: float64(parser.httpUnknown) / duration, Units: units, Timestamp: now},
	}, nil
}

func main() {}

var Parser SampleParser
