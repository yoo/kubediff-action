package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

}

type stringSlice []string

func (s *stringSlice) UnmarshalText(text []byte) error {
	lines := strings.Split(string(text), "\n")
	ret := []string{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		ret = append(ret, line)
	}
	*s = ret
	return nil
}

func run() error {
	config := &struct {
		GithubToken      string      `split_words:"true"`
		GithubRepository string      `required:"true" split_words:"true"`
		GithubHeadRef    string      `required:"true" split_words:"true"`
		CommentPR        bool        `envconfig:"INPUT_COMMENT_PR"`
		FilteredFields   stringSlice `envconfig:"INPUT_FILTERED_FIELDS"`
	}{}

	if err := envconfig.Process("", config); err != nil {
		return errors.Wrapf(err, "failed to load environment variables")
	}

	if len(os.Args) < 3 {
		return errors.New("two base directories required")
	}

	live := os.Args[1]
	merged := os.Args[2]

	Debugf("load diff files %q %q", live, merged)
	diffs, err := NewFileDiffs(live, merged)
	if err != nil {
		return errors.Wrapf(err, "failed to read manifests paths %q %q", live, merged)
	}

	if err = FilterFields(diffs, config.FilteredFields); err != nil {
		return errors.Wrap(err, "failed to filter fields from manifests")
	}

	if err = DiffFiles(diffs); err != nil {
		return errors.Wrapf(err, "failed to diff manifests from %q and %q", live, merged)
	}

	md, err := RenderMarkdown(diffs)
	if err != nil {
		return errors.Wrapf(err, "failed to render comment template")
	}

	fmt.Println(md)

	if !config.CommentPR {
		fmt.Println("\n pull request comment disabled")
		return nil
	}

	if config.GithubToken == "" {
		return errors.New("comment_pr set to \"true\", GITHUB_TOKEN required")
	}

	gh := NewGitHub(config.GithubToken)
	prID, err := gh.GetPullRequestID(config.GithubRepository, config.GithubHeadRef)
	if err != nil {
		return errors.Wrap(err, "failed to get pull request id")
	}
	err = gh.PrComment(config.GithubRepository, prID, md)
	return errors.Wrap(err, "failed to create pr comment")
}
