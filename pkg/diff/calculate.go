package diff

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"slices"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

func Calculate(
	ctx context.Context,
	repoPath string,
	commitA, commitB string,
) ([]*Diff, error) {
	cmd := exec.CommandContext(
		ctx,
		"git",
		"diff",
		"--patch-with-raw",
		commitA,
		commitB,
	)

	cmd.Dir = repoPath
	buffer := new(bytes.Buffer)
	cmd.Stdout = buffer

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if buffer.Len() == 0 {
		return nil, nil
	}

	diffs := make(map[string]*Diff)

	reader := bufio.NewReader(buffer)

	files, preamble, err := gitdiff.Parse(reader)
	if err != nil {
		return nil, err
	}

	preamble = strings.TrimSpace(preamble)
	for _, line := range strings.Split(preamble, "\n") {
		diff, err := parseDiff([]byte(line))
		if err != nil {
			return nil, err
		}

		diffs[diff.Src.Path] = &diff
	}

	for _, file := range files {
		var name string
		if file.OldName != "" {
			name = file.OldName
		} else {
			name = file.NewName
		}

		var lines []Line
		for _, fragment := range file.TextFragments {
			var (
				srcLineNumber = fragment.OldPosition
				dstLineNumber = fragment.NewPosition
			)

			for _, fragmentLine := range fragment.Lines {
				line := Line{
					Op:      fragmentLine.Op,
					Content: fragmentLine.Line,
				}

				switch fragmentLine.Op {
				case gitdiff.OpAdd:
					line.NumberInDst = dstLineNumber
					dstLineNumber++
				case gitdiff.OpDelete:
					line.NumberInSrc = srcLineNumber
					srcLineNumber++
				default:
					line.NumberInDst = dstLineNumber
					line.NumberInSrc = srcLineNumber
					dstLineNumber++
					srcLineNumber++
				}

				lines = append(lines, line)
			}
		}

		diff := diffs[name]
		diff.Lines = lines
		diff.IsBinary = file.IsBinary

		diffs[name] = diff
	}

	diffsSlice := make([]*Diff, 0, len(diffs))
	for _, diff := range diffs {
		diffsSlice = append(diffsSlice, diff)
	}

	slices.SortFunc(diffsSlice, func(a, b *Diff) int {
		return strings.Compare(a.Src.Path, b.Src.Path)
	})

	return diffsSlice, nil
}
