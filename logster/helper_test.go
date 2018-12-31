package logster

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadParserPlugin(t *testing.T) {
	pluginPath := os.Getenv("UT_PARSER_PLUGIN_PATH")
	_, err := LoadParserPlugin(pluginPath)
	assert.Nil(t, err)
}

func TestLoadOutputPlugin(t *testing.T) {
	pluginPath := os.Getenv("UT_OUTPUT_PLUGIN_PATH")
	_, err := LoadOutputPlugin(pluginPath)
	assert.Nil(t, err)
}

func TestLoadParserPluginError(t *testing.T) {
	pluginPath := os.Getenv("UT_OUTPUT_PLUGIN_PATH")
	_, err := LoadParserPlugin(pluginPath)
	assert.NotNil(t, err)

	invalidPath := filepath.Join("dummy", pluginPath)
	_, err = LoadParserPlugin(invalidPath)
	assert.NotNil(t, err)
}

func TestLoadOutputPluginError(t *testing.T) {
	pluginPath := os.Getenv("UT_PARSER_PLUGIN_PATH")
	_, err := LoadOutputPlugin(pluginPath)
	assert.NotNil(t, err)

	invalidPath := filepath.Join("dummy", pluginPath)
	_, err = LoadOutputPlugin(invalidPath)
	assert.NotNil(t, err)
}
