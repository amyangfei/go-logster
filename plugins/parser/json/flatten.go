package main

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Refer to: https://github.com/jeremywohl/flatten

// Flatten generates a flat map from a nested one.
// The original may include values of type map, slice and scalar, but not struct.
// Keys in the flat map will be a compound of descending map keys and slice iterations.
// keys are separated with separator. A prefix is joined to each key.
func Flatten(nested map[string]interface{}, prefix, separator string) (map[string]interface{}, error) {
	flatmap := make(map[string]interface{})

	err := flatten(true, flatmap, nested, prefix, separator)
	if err != nil {
		return nil, err
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
		return "", err
	}

	flatmap, err := Flatten(nested, prefix, separator)
	if err != nil {
		return "", err
	}

	flatb, err := json.Marshal(&flatmap)
	if err != nil {
		return "", err
	}

	return string(flatb), nil
}

func flatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix, separator string) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}, []interface{}:
			if err := flatten(false, flatMap, v, newKey, separator); err != nil {
				return err
			}
		default:
			flatMap[newKey] = v
		}

		return nil
	}

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
		return errors.New("Not a valid input: map or slice")
	}

	return nil
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
