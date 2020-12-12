package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type GitHub struct {
	ctx    context.Context
	client *github.Client
}

func NewGitHub(token string) *GitHub {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &GitHub{
		ctx:    ctx,
		client: github.NewClient(tc),
	}
}

func (g *GitHub) GetPullRequestID(repository, headRef string) (int, error) {
	ghr := strings.SplitN(repository, "/", 2)
	owner := ghr[0]
	repo := ghr[1]

	opts := &github.PullRequestListOptions{
		State: "open",
		Head:  headRef,
	}
	prs, _, err := g.client.PullRequests.List(g.ctx, owner, repo, opts)
	if err != nil {
		return 0, errors.Wrapf(err, "list pull requests for %q", repository)
	}

	if len(prs) == 0 {
		return 0, errors.Errorf("head ref %q not found in %q", headRef, repository)
	}
	return *prs[0].Number, nil
}
func (g *GitHub) PrComment(repository string, id int, text string) error {
	ghr := strings.SplitN(repository, "/", 2)
	owner := ghr[0]
	repo := ghr[1]

	newComment := &github.IssueComment{
		Body: &text,
	}

	oldComment, ok, err := g.findComment(owner, repo, id)
	if err != nil {
		return err
	}
	if ok {
		_, _, err = g.client.Issues.EditComment(g.ctx, owner, repo, *oldComment.ID, newComment)
		return errors.Wrapf(err, "edit comment in %q on pull request %q", repository, id)
	}

	_, _, err = g.client.Issues.CreateComment(g.ctx, owner, repo, id, newComment)
	return errors.Wrapf(err, "create comment in %q on pull request %d", repository, id)
}

func (g *GitHub) findComment(owner, repo string, id int) (*github.IssueComment, bool, error) {
	comments, _, err := g.client.Issues.ListComments(g.ctx, owner, repo, id, nil)
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to list comments in \"%s/%s\" on pull request %q", owner, repo, id)
	}

	currentUser, err := g.currentUser()
	if err != nil {
		return nil, false, errors.Wrap(err, "get current user")
	}

	for _, comment := range comments {
		if comment.User.GetLogin() != currentUser {
			continue
		}
		if strings.HasPrefix(*comment.Body, "### KubeDiff Action") {
			return comment, true, nil
		}
	}
	return nil, false, nil
}

func (g *GitHub) currentUser() (string, error) {
	user, resp, err := g.client.Users.Get(g.ctx, "")
	if resp.StatusCode == http.StatusForbidden {
		return "github-actions[bot]", nil
	}
	if err != nil {
		return "", errors.Wrap(err, "get current user")
	}
	return user.GetLogin(), nil
}
