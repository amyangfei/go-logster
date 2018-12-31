package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/amyangfei/go-logster/inter"
)

func CaptureOutput(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	f()
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

type TestData struct {
	name  string
	value float64
	ts    int64
}

func prepareData() []TestData {
	inputData := []TestData{
		{"name1", 1.1, time.Now().Unix()},
		{"name2", 2.1, time.Now().Unix()},
		{"name3", 3.1, time.Now().Unix()},
	}
	return inputData
}

var Logger zerolog.Logger

func init() {
	Logger = zerolog.New(os.Stdout)
}

func TestStdoutOutput(t *testing.T) {
	output := &StdoutOutput{}
	err := output.Init("P", "S", `{"separator": "-"}`, false, Logger)
	assert.Nil(t, err)

	inputData := prepareData()
	metrics := []*inter.Metric{}
	expected := ""
	for _, line := range inputData {
		metrics = append(metrics,
			&inter.Metric{Name: line.name, Value: line.value, Timestamp: line.ts})
		expected += fmt.Sprintf("%d P-%s-S %v\n", line.ts, line.name, line.value)
	}
	capturedOutput := CaptureOutput(func() { output.Submit(metrics) })
	assert.Equal(t, expected, capturedOutput)
}

func TestStdoutWithNoOption(t *testing.T) {
	output := &StdoutOutput{}
	err := output.Init("P", "", "", false, Logger)
	assert.Nil(t, err)
	inputData := prepareData()
	metrics := []*inter.Metric{}
	expected := ""
	for _, line := range inputData {
		metrics = append(metrics,
			&inter.Metric{Name: line.name, Value: line.value, Timestamp: line.ts})
		expected += fmt.Sprintf("%d P.%s %v\n", line.ts, line.name, line.value)
	}
	capturedOutput := CaptureOutput(func() { output.Submit(metrics) })
	assert.Equal(t, expected, capturedOutput)
}
