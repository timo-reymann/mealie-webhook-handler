package mealie

import (
	"fmt"
	"io"
	"net/http"
)

func FetchRecipeImage(apiUrl string, recipeId string, version string) ([]byte, error) {
	res, err := http.Get(fmt.Sprintf("%s/media/recipes/%s/images/original.webp?version=%s", apiUrl, recipeId, version))
	if err != nil {
		return nil, err
	}

	if res.Header.Get("Content-Type") == "application/json" {
		return nil, nil
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
