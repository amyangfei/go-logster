package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getJSONParser(options string) *JSONParser {
	p := &JSONParser{}
	p.Init(options)
	return p
}

func TestValidJson(t *testing.T) {
	p := getJSONParser(`{"separator": "$", "prefix": "T"}`)

	line := `{"1.1":
				{"value1": 0,
			 	 "value2": "hi",
				 "1.2": {"value3": 0.1, "value4": false}},
				 "2.1": ["a", "b"]
			 }`
	expected := map[string]interface{}{
		"T1.1$value1":     0,
		"T1.1$value2":     "hi",
		"T1.1$1.2$value3": 0.1,
		"T1.1$1.2$value4": false,
		"T2.1$0":          "a",
		"T2.1$1":          "b",
	}

	err := p.ParseLine(line)
	assert.Nil(t, err)

	metrics, err := p.GetState(0)
	assert.Nil(t, err)
	for _, metric := range metrics {
		assert.Contains(t, expected, metric.Name)
		assert.EqualValues(t, expected[metric.Name], metric.Value)
	}
}

func TestJsonMerge(t *testing.T) {
	p := getJSONParser(`{"separator": "."}`)

	lines := []string{
		`{"1.1": {
			"value1": 0,
			"value2": "hi",
			"1.2": {"value3": 0.1, "value4": false}
		  },
		  "2.1": ["a", "b"],
		  "3.1": "first"
		 }`,
		`{"2.1": ["c"],
		  "3.1": "second"
		 }`,
		`{"2.1": ["d", -1],
		  "3.1": "third"
		 }`,
	}
	expected := map[string]interface{}{
		"1.1.value1":     0,
		"1.1.value2":     "hi",
		"1.1.1.2.value3": 0.1,
		"1.1.1.2.value4": false,
		"2.1.0":          "a",
		"2.1.1":          "b",
		"2.1.2":          "c",
		"2.1.3":          "d",
		"2.1.4":          -1,
		"3.1":            "first",
	}

	for _, line := range lines {
		err := p.ParseLine(line)
		assert.Nil(t, err)
	}

	metrics, err := p.GetState(0)
	assert.Nil(t, err)
	for _, metric := range metrics {
		assert.Contains(t, expected, metric.Name)
		assert.EqualValues(t, expected[metric.Name], metric.Value)
	}
}

func TestInvalidJson(t *testing.T) {
	p := getJSONParser("")
	err := p.ParseLine(`{"hello": "world"`)
	assert.NotNil(t, err)
}
