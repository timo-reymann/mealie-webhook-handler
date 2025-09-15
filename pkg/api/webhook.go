package api

import (
	"io"
	"net/http"

	"go.deepl.dev/mealie-webhook-handler/pkg/appcontext"
	"go.deepl.dev/mealie-webhook-handler/pkg/configuration"
	"go.deepl.dev/mealie-webhook-handler/pkg/mealie"
	"go.deepl.dev/mealie-webhook-handler/pkg/output"
	"go.deepl.dev/mealie-webhook-handler/pkg/template"
)

func CreateHandleWebhook(appCtx appcontext.AppContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			SendError(w, http.StatusMethodNotAllowed, "Only POST is allowed")
			return
		}

		hookId := r.PathValue("identifier")
		webhook, ok := appCtx.Config.Webhooks[hookId]

		if !ok {
			SendError(w, http.StatusNotFound, "No webhook registered for identifier")
			return
		}

		rawPayload, err := io.ReadAll(r.Body)
		if err != nil {
			SendError(w, http.StatusInternalServerError, "No body")
			return
		}

		parsedPayload, err := mealie.ParseWebhook(rawPayload)
		if err != nil {
			SendError(w, http.StatusInternalServerError, "Failed to parse payload: "+err.Error())
			return
		}

		image, _ := mealie.FetchRecipeImage(appCtx.Config.Mealie.ApiUrl, parsedPayload.Content.ID, parsedPayload.Content.Image)

		tplPayload := configuration.OutputConfigTemplatePayload{
			Recipe:   parsedPayload.Content,
			Servings: parsedPayload.RecipeScale,
			HasImage: image != nil,
		}

		options, err := webhook.TemplateOptions(tplPayload)
		if err != nil {
			SendError(w, http.StatusPreconditionFailed, "Invalid template")
			return
		}

		out := output.Outputs()[webhook.Output]
		recipeTpl, err := webhook.LoadRecipeTemplate()
		if err != nil {
			SendError(w, http.StatusInternalServerError, "template file not found")
			return
		}

		templatedRecipe, err := template.Exec("recipe", string(recipeTpl), tplPayload)
		if err != nil {
			SendError(w, http.StatusInternalServerError, "recipe could not be templated: "+err.Error())
			return
		}

		err = out.Output(r.Context(), *templatedRecipe, image, options)
		if err != nil {
			SendError(w, http.StatusInternalServerError, "failed to run output: "+err.Error())
		}

		_, _ = w.Write([]byte(hookId))
	}
}
