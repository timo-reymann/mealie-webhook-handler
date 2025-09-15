package configuration

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"go.deepl.dev/mealie-webhook-handler/pkg/output"
)

func ParseConfiguration(config []byte) (*Config, error) {
	var model Config
	err := toml.Unmarshal(config, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *Config) Init() error {
	outs := output.Outputs()
	for id, webhook := range c.Webhooks {
		outputName := webhook.Output

		out, ok := outs[outputName]
		if !ok {
			return fmt.Errorf("webhook %s has no valid output", id)
		}

		err := out.Init()
		if err != nil {
			return fmt.Errorf("webhook %s could not be initialized as output init faile with %s", id, err)
		}

		err = out.ValidateOptions(webhook.OutputOptions)
		if err != nil {
			return fmt.Errorf("webhook %s has invalid output options %s", id, err)
		}
	}
	return nil
}
