package main

import (
	"encoding/json"
	"strconv"

	"github.com/juju/errors"
)

// Refer to: https://github.com/jeremywohl/flatten

var (
	// ErrorInvalidInput is generic error for invalid json input
	ErrorInvalidInput = errors.New("Not a valid input: map or slice")
)

// Flatten generates a flat map from a nested one.
// The original may include values of type map, slice and scalar, but not struct.
// Keys in the flat map will be a compound of descending map keys and slice iterations.
// keys are separated with separator. A prefix is joined to each key.
func Flatten(nested map[string]interface{}, prefix, separator string) (map[string]interface{}, error) {
	flatmap := make(map[string]interface{})

	err := flatten(true, flatmap, nested, prefix, separator)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return flatmap, nil
}

// FlattenString generates a flat JSON map from a nested one.
// Keys in the flat map will be a compound of descending map keys and slice iterations.
// keys are separated with separator. A prefix is joined to each key.
func FlattenString(nestedstr, prefix, separator string) (string, error) {
	var nested map[string]interface{}
	err := json.Unmarshal([]byte(nestedstr), &nested)
	if err != nil {
		return "", errors.Trace(err)
	}

	flatmap, err := Flatten(nested, prefix, separator)
	if err != nil {
		return "", errors.Trace(err)
	}

	flatb, err := json.Marshal(&flatmap)
	if err != nil {
		return "", errors.Trace(err)
	}

	return string(flatb), nil
}

func flatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix, separator string) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}, []interface{}:
			// gofail: var FlattenError2 bool
			// if FlattenError2 {
			//   v = map[int]int{1:2}
			// }
			// return errors.Trace(flatten(false, flatMap, v, newKey, separator))
			return errors.Trace(flatten(false, flatMap, v, newKey, separator))
		default:
			flatMap[newKey] = v
		}

		return nil
	}

	dispatch := func(obj interface{}) error {
		switch nested.(type) {
		case map[string]interface{}:
			for k, v := range nested.(map[string]interface{}) {
				newKey := enkey(top, prefix, k, separator)
				assign(newKey, v)
			}
		case []interface{}:
			for i, v := range nested.([]interface{}) {
				newKey := enkey(top, prefix, strconv.Itoa(i), separator)
				assign(newKey, v)
			}
		default:
			return ErrorInvalidInput
		}
		return nil
	}

	// gofail: var FlattenError1 bool
	// if FlattenError1 {
	//   nested = map[int]int{1:2}
	// }
	// return errors.Trace(dispatch(nested))
	return errors.Trace(dispatch(nested))
}

func enkey(top bool, prefix, subkey, separator string) string {
	key := prefix

	if top {
		key += subkey
	} else {
		key += separator + subkey
	}

	return key
}
