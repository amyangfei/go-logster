package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"unicode"

	gofail "github.com/etcd-io/gofail/runtime"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	cases := []struct {
		test      string
		want      map[string]interface{}
		prefix    string
		separator string
	}{
		{
			`{
                "foo": {
                    "jim":"bean"
                },
                "fee": "bar",
                "n1": {
                    "alist": [
                        "a",
                        "b",
                        "c",
                        {
                            "d": "other",
                            "e": "another"
                        }
                    ]
                },
                "number": 1.4567,
                "bool":   true
            }`,
			map[string]interface{}{
				"foo.jim":      "bean",
				"fee":          "bar",
				"n1.alist.0":   "a",
				"n1.alist.1":   "b",
				"n1.alist.2":   "c",
				"n1.alist.3.d": "other",
				"n1.alist.3.e": "another",
				"number":       1.4567,
				"bool":         true,
			},
			"",
			".",
		}, {
			`{ "a": { "b": "c" }, "e": "f" }`,
			map[string]interface{}{
				"p:a.b": "c",
				"p:e":   "f",
			},
			"p:",
			".",
		},
	}

	for i, test := range cases {
		var m interface{}
		err := json.Unmarshal([]byte(test.test), &m)
		if err != nil {
			t.Errorf("%d: failed to unmarshal test: %v", i+1, err)
		}
		got, err := Flatten(m.(map[string]interface{}), test.prefix, test.separator)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

func TestFlattenString(t *testing.T) {
	cases := []struct {
		test      string
		want      string
		prefix    string
		separator string
	}{
		{
			`{ "a": "b" }`,
			`{ "a": "b" }`,
			"",
			".",
		},
		{
			`{ "a": { "b" : { "c" : { "d" : "e" } } }, "number": 1.4567, "bool": true }`,
			`{ "a.b.c.d": "e", "bool": true, "number": 1.4567 }`,
			"",
			".",
		},
	}

	for i, test := range cases {
		got, err := FlattenString(test.test, test.prefix, test.separator)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}

		nixws := func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}

		if got != strings.Map(nixws, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

func TestFlattenError(t *testing.T) {
	prefix := ""
	separator := "."
	_, err := FlattenString(`{"a": "b", "c": }`, prefix, separator)
	assert.NotNil(t, err)

	defer gofail.Enable("github.com/amyangfei/go-logster/plugins/parser/json/FlattenError1", `return(false)`)
	defer gofail.Enable("github.com/amyangfei/go-logster/plugins/parser/json/FlattenError2", `return(false)`)

	gofail.Enable("github.com/amyangfei/go-logster/plugins/parser/json/FlattenError1", `return(true)`)
	_, err = FlattenString(`{"a": "b"}`, prefix, separator)
	assert.Equal(t, ErrorInvalidInput, errors.Cause(err))

	gofail.Enable("github.com/amyangfei/go-logster/plugins/parser/json/FlattenError2", `return(true)`)
	_, err = FlattenString(`{"a": "b"}`, prefix, separator)
	assert.Equal(t, ErrorInvalidInput, errors.Cause(err))
}
