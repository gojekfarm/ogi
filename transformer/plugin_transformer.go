package ogitransformer

import (
	"fmt"
	"path"
	"plugin"

	"github.com/abhishekkr/gol/golenv"

	logger "github.com/gojektech/ogi/logger"
)

type TransformerPlugin struct {
	Name          string
	TransformFunc plugin.Symbol
}

var (
	TransformerPluginPath = golenv.OverrideIfEnv("TRANSFORMER_PLUGIN_PATH", "./transformer.so")
)

func NewTransformerPlugin() Transformer {
	fmt.Println(TransformerPluginPath)
	p, err := plugin.Open(TransformerPluginPath)
	if err != nil {
		logger.Fatalf("Error reading plugin: %s", err)
	}

	transformFunc, err := p.Lookup("Transform")
	if err != nil {
		logger.Fatalf("Error looking up 'Transform': %s", err)
	}

	return &TransformerPlugin{
		Name:          path.Base(TransformerPluginPath),
		TransformFunc: transformFunc,
	}
}

func (plugin *TransformerPlugin) Transform(msg []byte) error {
	return plugin.TransformFunc.(func([]byte) error)(msg)
}
