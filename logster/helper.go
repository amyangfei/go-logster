package logster

import (
	"plugin"

	"github.com/juju/errors"

	"github.com/amyangfei/go-logster/inter"
)

// LoadParserPlugin loads a parser plugin from given plugin path,
// the plugin name is specificed to Parser
func LoadParserPlugin(pluginPath string) (inter.Parser, error) {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	symbol, err := plug.Lookup("Parser")
	if err != nil {
		return nil, err
	}
	parser, ok := symbol.(inter.Parser)
	if !ok {
		return nil, errors.New("unexpected type from module symbol")
	}
	return parser, nil
}

// LoadOutputPlugin loads a Output plugin from given plugin path,
// the plugin name is specificed to Output
func LoadOutputPlugin(pluginPath string) (inter.Output, error) {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	symbol, err := plug.Lookup("Output")
	if err != nil {
		return nil, err
	}
	output, ok := symbol.(inter.Output)
	if !ok {
		return nil, errors.New("unexpected type from module symbol")
	}
	return output, nil
}
