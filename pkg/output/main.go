package output

import (
	"context"

	"go.deepl.dev/mealie-webhook-handler/pkg/output/github_pr"
)

type Output interface {
	Name() string
	Init() error
	ValidateOptions(options map[string]string) error
	Output(ctx context.Context, templatedRecipe string, image []byte, config map[string]string) error
}

var outputs = map[string]Output{}

// Register given output
func Register(output Output) {
	name := output.Name()
	outputs[name] = output
}

func Outputs() map[string]Output {
	return outputs
}

func init() {
	Register(&github_pr.GitHubPullRequestOutput{})
}
