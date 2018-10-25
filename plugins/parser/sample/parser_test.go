package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getSampleParser() *SampleParser {
	p := &SampleParser{}
	p.Init("")
	return p
}

func preParseLog(p *SampleParser) []error {
	errors := make([]error, 0)
	accessLogTmpl := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" %s 2326`
	errors = append(errors, p.ParseLine(fmt.Sprintf(accessLogTmpl, "200")))
	errors = append(errors, p.ParseLine(fmt.Sprintf(accessLogTmpl, "201")))
	errors = append(errors, p.ParseLine(fmt.Sprintf(accessLogTmpl, "400")))
	return errors
}

func TestValidLines(t *testing.T) {
	p := getSampleParser()
	errors := preParseLog(p)
	for _, err := range errors {
		assert.Nil(t, err)
	}

	assert.Equal(t, p.http1xx, 0)
	assert.Equal(t, p.http2xx, 2)
	assert.Equal(t, p.http3xx, 0)
	assert.Equal(t, p.http4xx, 1)
	assert.Equal(t, p.http5xx, 0)
	assert.Equal(t, p.httpUnknown, 0)
}

func TestMetrics1Sec(t *testing.T) {
	p := getSampleParser()
	preParseLog(p)
	metrics, err := p.GetState(1)
	assert.Nil(t, err)
	expected := map[string]float64{
		"http_1xx":     0,
		"http_2xx":     2,
		"http_3xx":     0,
		"http_4xx":     1,
		"http_5xx":     0,
		"http_unknown": 0,
	}
	for _, metric := range metrics {
		assert.Contains(t, expected, metric.Name)
		assert.Equal(t, expected[metric.Name], metric.Value)
	}
}

func TestMetrics2Sec(t *testing.T) {
	p := getSampleParser()
	preParseLog(p)
	metrics, err := p.GetState(2)
	assert.Nil(t, err)
	expected := map[string]interface{}{
		"http_1xx":     0,
		"http_2xx":     1,
		"http_3xx":     0,
		"http_4xx":     0.5,
		"http_5xx":     0,
		"http_unknown": 0,
	}
	for _, metric := range metrics {
		assert.Contains(t, expected, metric.Name)
		assert.EqualValues(t, expected[metric.Name], metric.Value)
	}
}

func TestInvalidLine(t *testing.T) {
	p := getSampleParser()
	err := p.ParseLine("invalid log entry")
	assert.NotNil(t, err)
}
