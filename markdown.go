package main

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

const commentTemplate = "### KubeDiff Action\n{{range .}}Object `{{.ObjectPath}}`:\n```diff\n{{.Diff}}\n```\n{{end}}"

func RenderMarkdown(diffs []*FileDiff) (string, error) {
	t := template.New("comment")
	t, err := t.Parse(commentTemplate)
	if err != nil {
		return "", errors.Wrap(err, "parse template")
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, diffs)
	return buf.String(), errors.Wrap(err, "render template")
}
