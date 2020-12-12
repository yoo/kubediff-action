package main

import (
	"bytes"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type FileDiff struct {
	ObjectPath string
	LivePath   string
	LiveBuf    []byte
	MergedPath string
	MergedBuf  []byte
	Diff       string
}

func NewFileDiffs(livePath, mergedPath string) ([]*FileDiff, error) {
	fileInfos, err := ioutil.ReadDir(livePath)
	if err != nil {
		return nil, errors.Wrapf(err, "read dir %q", livePath)
	}

	files := []string{}
	for _, fi := range fileInfos {
		if fi.IsDir() {
			continue
		}
		files = append(files, fi.Name())
	}
	sort.Strings(files)
	Debugf("found yaml files %v", files)

	diffs := []*FileDiff{}
	for _, file := range files {
		liveBuf, err := ioutil.ReadFile(path.Join(livePath, file))
		if err != nil {
			return nil, errors.Wrapf(err, "read file %q", livePath)
		}
		mergedBuf, err := ioutil.ReadFile(path.Join(mergedPath, file))
		if err != nil {
			return nil, errors.Wrapf(err, "read file %q", mergedPath)
		}
		diff := &FileDiff{
			ObjectPath: file,
			LivePath:   path.Join(livePath, file),
			LiveBuf:    liveBuf,
			MergedPath: path.Join(mergedPath, file),
			MergedBuf:  mergedBuf,
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func DiffFiles(diffs []*FileDiff) error {
	for _, diff := range diffs {
		if err := diff.doDiff(); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileDiff) doDiff() error {
	dmp := diffmatchpatch.New()
	liveDMP, mergedDMP, dmpStrings := dmp.DiffLinesToChars(string(f.LiveBuf), string(f.MergedBuf))
	diffs := dmp.DiffMain(liveDMP, mergedDMP, true)
	diffs = dmp.DiffCharsToLines(diffs, dmpStrings)
	diffs = dmp.DiffCleanupSemantic(diffs)
	f.Diff = f.diffsToString(diffs)
	return nil
}

func (f *FileDiff) diffsToString(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := strings.Trim(diff.Text, "\n")
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buff.WriteString("\n+ ")
			for _, r := range text {
				if r == '\n' {
					buff.WriteString("\n+ ")
					continue
				}
				buff.WriteRune(r)
			}
		case diffmatchpatch.DiffDelete:
			buff.WriteString("\n- ")
			for _, r := range text {
				if r == '\n' {
					buff.WriteString("\n- ")
					continue
				}
				buff.WriteRune(r)
			}
		case diffmatchpatch.DiffEqual:
			buff.WriteString("\n  ")
			for _, r := range text {
				if r == '\n' {
					buff.WriteString("\n  ")
					continue
				}
				buff.WriteRune(r)
			}
		}
	}
	return strings.Trim(buff.String(), "\n")

}
