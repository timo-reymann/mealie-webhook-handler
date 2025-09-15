package mealie

import "encoding/json"

func ParseWebhook(raw []byte) (Recipe, error) {
	var recipe Recipe
	err := json.Unmarshal(raw, &recipe)
	return recipe, err
}
