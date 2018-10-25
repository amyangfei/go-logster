package main

import (
	"encoding/json"
	"time"

	"github.com/amyangfei/go-logster/logster"
	"github.com/buger/jsonparser"
	"github.com/imdario/mergo"
)

const DefaultKeySeparator = "."
const DefaultPrefix = ""

type JsonParser struct {
	KeySeparator string
	Prefix       string
	Metrics      map[string]interface{}
}

func parserKey(options, key, defaultVal string) (string, error) {
	val, dataType, _, err := jsonparser.Get([]byte(options), key)
	if err != nil {
		return "", err
	}
	var result string
	if dataType == jsonparser.NotExist {
		result = defaultVal
	} else {
		result = string(val)
	}
	return result, nil
}

func (parser *JsonParser) Init(options string) error {
	if val, err := parserKey(options, "seperator", DefaultKeySeparator); err != nil {
		return err
	} else {
		parser.KeySeparator = val
	}
	if val, err := parserKey(options, "prefix", DefaultPrefix); err != nil {
		return err
	} else {
		parser.Prefix = val
	}
	parser.Metrics = make(map[string]interface{})
	return nil
}

func (parser *JsonParser) ParseLine(line string) error {
	var nested map[string]interface{}
	if err := json.Unmarshal([]byte(line), &nested); err != nil {
		return err
	}
	mergo.Merge(&parser.Metrics, nested, mergo.WithAppendSlice)
	return nil
}

func (parser *JsonParser) GetState(duration float64) ([]*logster.Metric, error) {
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

var Parser JsonParser
