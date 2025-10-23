package main

import (
	"syscall/js"

	"go.deepl.dev/mealie-webhook-handler/pkg/configuration"
	"go.deepl.dev/mealie-webhook-handler/pkg/mealie"
	"go.deepl.dev/mealie-webhook-handler/pkg/template"
)

func main() {
	js.Global().
		Set("renderWebhookTemplate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) != 2 {
				return "Invalid number of arguments"
			}

			payload := args[0].String()
			parsed, err := mealie.ParseWebhook([]byte(payload))
			if err != nil {
				return "Failed to parse payload: " + err.Error()
			}

			tpl := args[1].String()
			tplPayload := configuration.OutputConfigTemplatePayload{
				Recipe:   parsed.Content,
				Servings: parsed.RecipeScale,
				HasImage: true,
			}
			templatedRecipe, err := template.Exec("recipe", tpl, tplPayload)
			return *templatedRecipe
		}))
	<-make(chan bool)
}
