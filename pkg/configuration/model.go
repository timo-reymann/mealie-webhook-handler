package configuration

import (
	"os"

	"go.deepl.dev/mealie-webhook-handler/pkg/mealie"
	"go.deepl.dev/mealie-webhook-handler/pkg/template"
)

type WebhookConfig struct {
	TemplatePath  string            `toml:"template_path"`
	Output        string            `toml:"output"`
	OutputOptions map[string]string `toml:"output_options"`
}

type OutputConfigTemplatePayload struct {
	Recipe   mealie.RecipeContent
	Servings float64
	HasImage bool
}

func (wc *WebhookConfig) TemplateOptions(payload OutputConfigTemplatePayload) (map[string]string, error) {
	options := map[string]string{}
	for k, v := range wc.OutputOptions {
		rendered, err := template.Exec(k, v, payload)
		if err != nil {
			return nil, err
		}
		options[k] = *rendered
	}
	return options, nil
}

func (wc *WebhookConfig) LoadRecipeTemplate() ([]byte, error) {
	return os.ReadFile(wc.TemplatePath)
}

type Mealie struct {
	ApiUrl string `toml:"api_url"`
}

type Config struct {
	Webhooks map[string]WebhookConfig `toml:"webhook"`
	Mealie   Mealie                   `toml:"mealie"`
}
