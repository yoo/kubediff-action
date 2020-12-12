package main

import (
	"bytes"

	"github.com/gonvenience/ytbx"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func FilterFields(diffs []*FileDiff, fields []string) error {
	for _, diff := range diffs {
		Debugf("filter %q and %q", diff.LivePath, diff.MergedPath)
		live, merged, err := ytbx.LoadFiles(diff.LivePath, diff.MergedPath)
		if err != nil {
			return errors.Wrapf(err, "load yaml file")
		}
		diff.LiveBuf, err = filterFields(live.Documents, fields)
		if err != nil {
			return errors.Wrapf(err, "filter %q", diff.LivePath)
		}
		diff.MergedBuf, err = filterFields(merged.Documents, fields)
		if err != nil {
			return errors.Wrapf(err, "filter %q", diff.MergedPath)
		}
	}
	return nil
}

func filterFields(documents []*yaml.Node, fields []string) ([]byte, error) {
	buf := &bytes.Buffer{}
	for _, document := range documents {
		for _, field := range fields {
			Debugf("filter field %q", field)
			// nolint:errcheck
			ytbx.Delete(document, field)
		}
		err := yaml.NewEncoder(buf).Encode(document)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to yaml encode after filtering")
		}
	}
	return buf.Bytes(), nil
}
