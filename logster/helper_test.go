package logster

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadParserPlugin(t *testing.T) {
	pluginPath := os.Getenv("UT_PARSER_PLUGIN_PATH")
	t.Log(pluginPath)
	_, err := LoadParserPlugin(pluginPath)
	assert.Nil(t, err)
}

func TestLoadOutputPlugin(t *testing.T) {
	pluginPath := os.Getenv("UT_OUTPUT_PLUGIN_PATH")
	t.Log(pluginPath)
	_, err := LoadOutputPlugin(pluginPath)
	assert.Nil(t, err)
}
