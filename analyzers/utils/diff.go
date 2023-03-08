package utils

import (
	"bytes"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/tools/go/analysis"
)

var hunkRE = regexp.MustCompile(`@@ -(\d+),(\d+) \+\d+,\d+ @@`)

func GetSuggestedFixesFromDiff(file *token.File, a, b []byte) ([]analysis.SuggestedFix, error) {
	d := difflib.UnifiedDiff{
		A:       difflib.SplitLines(string(a)),
		B:       difflib.SplitLines(string(b)),
		Context: 1,
	}
	diff, err := difflib.GetUnifiedDiffString(d)
	if err != nil {
		return nil, err
	}
	if diff == "" {
		return nil, nil
	}
	var (
		fix   analysis.SuggestedFix
		first = true
		edit  analysis.TextEdit
		buf   bytes.Buffer
	)
	for _, line := range strings.Split(diff, "\n") {
		if line == "" {
			continue
		}
		hunk := hunkRE.FindStringSubmatch(line)
		switch {
		case len(hunk) > 0:
			if !first {
				edit.NewText = buf.Bytes()
				buf = bytes.Buffer{}
				fix.TextEdits = append(fix.TextEdits, edit)
				edit = analysis.TextEdit{}
			}
			first = false
			start, err := strconv.Atoi(hunk[1])
			if err != nil {
				return nil, err
			}
			lines, err := strconv.Atoi(hunk[2])
			if err != nil {
				return nil, err
			}
			edit.Pos = file.LineStart(start)
			end := start + lines
			if end > file.LineCount() {
				edit.End = token.Pos(file.Size())
			} else {
				edit.End = file.LineStart(end)
			}
		case line[0] == '+':
			buf.WriteString(line[1:])
			buf.WriteRune('\n')
		case line[0] == '-':
			// just skip
		default:
			buf.WriteString(line)
			buf.WriteRune('\n')
		}
	}
	edit.NewText = buf.Bytes()
	fix.TextEdits = append(fix.TextEdits, edit)

	return []analysis.SuggestedFix{fix}, nil
}
