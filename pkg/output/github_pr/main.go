package github_pr

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/google/go-github/v80/github"
	"go.deepl.dev/mealie-webhook-handler/pkg/output/validation"
)

const sourceBranchCfgKey = "source_branch"
const targetBranchCfgKey = "target_branch"
const repoSlugCfgKey = "slug"
const prTitleCfgKey = "title"
const commitMsgCfgKey = "commit_message"
const recipePathCfgKey = "recipe_path"
const imagePathCfgKey = "image_path"
const prBodyCfgKey = "body"

type GitHubPullRequestOutput struct {
	github *github.Client
}

func (g *GitHubPullRequestOutput) Init() error {
	if g.github == nil {
		g.github = github.NewClient(nil).
			WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	}
	return nil
}

func (g *GitHubPullRequestOutput) Name() string {
	return "github_pr"
}

func (g *GitHubPullRequestOutput) ValidateOptions(options map[string]string) error {
	return validation.FailOnFirst(
		validation.RequireKey(prTitleCfgKey),
		validation.RequireKey(prBodyCfgKey),
		validation.RequireKey(repoSlugCfgKey),
		validation.RequireKey(sourceBranchCfgKey),
		validation.RequireKey(targetBranchCfgKey),
		validation.RequireKey(recipePathCfgKey),
		validation.RequireKey(commitMsgCfgKey),
		validation.RequireKey(imagePathCfgKey),
		func(m map[string]string) error {
			slug, _ := m[repoSlugCfgKey]
			parts := strings.Split(slug, "/")
			if len(parts) != 2 {
				return errors.New("slug is malformed")
			}
			return nil
		},
	)(options)
}

func (g *GitHubPullRequestOutput) Output(ctx context.Context, templatedRecipe string, image []byte, config map[string]string) error {
	repoParts := strings.Split(config[repoSlugCfgKey], "/")
	owner := repoParts[0]
	repo := repoParts[1]

	targetBranch, _, err := g.github.Git.GetRef(ctx, owner, repo, fmt.Sprintf("heads/%s", config[targetBranchCfgKey]))
	if err != nil {
		return err
	}
	baseCommitSHA := targetBranch.Object.GetSHA()

	_, r, err := g.github.Git.CreateRef(
		ctx,
		owner,
		repo,
		github.CreateRef{
			Ref: fmt.Sprintf("refs/heads/%s", config[sourceBranchCfgKey]),
			SHA: baseCommitSHA,
		},
	)
	if err != nil && (r != nil && r.StatusCode != 422) {
		return err
	}

	recipeFileContent, _, _, err := g.github.Repositories.GetContents(ctx, owner, repo, config[recipePathCfgKey], &github.RepositoryContentGetOptions{Ref: config[sourceBranchCfgKey]})
	var recipefileSHA *string
	if err == nil && recipeFileContent != nil {
		recipefileSHA = recipeFileContent.SHA
	} else {
		recipefileSHA = nil
	}

	_, _, err = g.github.Repositories.CreateFile(
		ctx,
		owner,
		repo,
		config[recipePathCfgKey],
		&github.RepositoryContentFileOptions{
			Message: github.Ptr(config[commitMsgCfgKey]),
			Content: []byte(templatedRecipe),
			Branch:  github.Ptr(config[sourceBranchCfgKey]),
			SHA:     recipefileSHA,
		},
	)
	if err != nil {
		return err
	}

	if image != nil && len(image) > 1 {
		imageFileContent, _, _, err := g.github.Repositories.GetContents(ctx, owner, repo, config[imagePathCfgKey], &github.RepositoryContentGetOptions{Ref: config[sourceBranchCfgKey]})
		var imageFileSHA *string
		if err == nil && imageFileContent != nil {
			imageFileSHA = imageFileContent.SHA
		} else {
			imageFileSHA = nil
		}

		_, _, err = g.github.Repositories.CreateFile(
			ctx,
			owner,
			repo,
			config[imagePathCfgKey],
			&github.RepositoryContentFileOptions{
				Message: github.Ptr(config[commitMsgCfgKey]),
				Content: image,
				Branch:  github.Ptr(config[sourceBranchCfgKey]),
				SHA:     imageFileSHA,
			},
		)
		if err != nil {
			return err
		}
	}

	pr, r, err := g.github.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title:               github.Ptr(config[prTitleCfgKey]),
		Base:                github.Ptr(config[targetBranchCfgKey]),
		Head:                github.Ptr(config[sourceBranchCfgKey]),
		Body:                github.Ptr(config[prBodyCfgKey]),
		MaintainerCanModify: github.Ptr(true),
		Draft:               github.Ptr(false),
	})
	if err != nil && (r != nil && r.StatusCode != 422) {
		slog.Error("Failed to create PR", "err", err)
		return err
	}
	if pr != nil {

		slog.Info("Created PR", "id", pr.ID, "link", pr.Links.HTML)
	}
	return nil
}
