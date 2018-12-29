package main

import (
	"encoding/json"
	"time"

	"github.com/amyangfei/go-logster/logster"
	"github.com/buger/jsonparser"
	"github.com/imdario/mergo"
)

// DefaultKeySeparator is separator used in combined json key
const DefaultKeySeparator = "."

// DefaultPrefix is prefix added to json key
const DefaultPrefix = ""

// JSONParser holds a json parser
type JSONParser struct {
	KeySeparator string
	Prefix       string
	Metrics      map[string]interface{}
}

func parserKey(options, key, defaultVal string) (string, error) {
	val, dataType, _, err := jsonparser.Get([]byte(options), key)
	if err != nil {
		if dataType == jsonparser.NotExist {
			return defaultVal, nil
		}
		return "", err
	}
	return string(val), nil
}

// Init inits the *JSONParser type Parser
func (parser *JSONParser) Init(options string) error {
	val, err := parserKey(options, "separator", DefaultKeySeparator)
	if err != nil {
		return err
	}
	parser.KeySeparator = val

	val, err = parserKey(options, "prefix", DefaultPrefix)
	if err != nil {
		return err
	}
	parser.Prefix = val

	parser.Metrics = make(map[string]interface{})
	return nil
}

// ParseLine parses online and caches the parsed result
func (parser *JSONParser) ParseLine(line string) error {
	var nested map[string]interface{}
	if err := json.Unmarshal([]byte(line), &nested); err != nil {
		return err
	}
	mergo.Merge(&parser.Metrics, nested, mergo.WithAppendSlice)
	return nil
}

// GetState gets flatten json metrics from cached parsed result
func (parser *JSONParser) GetState(duration float64) ([]*logster.Metric, error) {
	flatmap, err := Flatten(parser.Metrics, parser.Prefix, parser.KeySeparator)
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	result := make([]*logster.Metric, 0)
	for k, v := range flatmap {
		result = append(result, &logster.Metric{Name: k, Value: v, Timestamp: now})
	}
	return result, nil
}

func main() {}

// Parser declares a JSONParser object
var Parser JSONParser
