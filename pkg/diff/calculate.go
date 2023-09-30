package diff

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/evermake/git-diff-view/pkg/gitutil"
)

func Calculate(
	ctx context.Context,
	repoPath string,
	commitA, commitB string,
) ([]*Diff, error) {
	exists, err := gitutil.RevisionExists(ctx, commitA)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("revision %s does not exists", commitA)
	}

	exists, err = gitutil.RevisionExists(ctx, commitB)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("revision %s does not exists", commitA)
	}

	cmd := exec.CommandContext(
		ctx,
		"git",
		"diff",
		"--patch-with-raw",
		commitA,
		commitB,
	)

	cmd.Dir = repoPath

	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// No changes
	if stdout.Len() == 0 {
		return nil, nil
	}

	diffs := make(map[string]*Diff)

	files, preamble, err := gitdiff.Parse(stdout)
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
				line := Line{}

				switch fragmentLine.Op {
				case gitdiff.OpAdd:
					line.Dst.Number = dstLineNumber
					line.Dst.Content = fragmentLine.Line

					line.Operation = LineOperationAdd

					dstLineNumber++
				case gitdiff.OpDelete:
					line.Src.Number = srcLineNumber
					line.Src.Content = fragmentLine.Line

					line.Operation = LineOperationDelete

					srcLineNumber++
				default:
					line.Dst.Number = dstLineNumber
					line.Dst.Content = fragmentLine.Line

					line.Src.Number = srcLineNumber
					line.Src.Content = fragmentLine.Line

					dstLineNumber++
					srcLineNumber++
				}

				lines = append(lines, line)
			}
		}

		var (
			deletedLinesIndices []int
			addedLinesIndices   []int
		)

		for i, line := range lines {
			switch line.Operation {
			case LineOperationAdd:
				addedLinesIndices = append(addedLinesIndices, i)
			case LineOperationDelete:
				deletedLinesIndices = append(deletedLinesIndices, i)
			default:
				window := min(len(deletedLinesIndices), len(addedLinesIndices))

				// update operation to modify for deleted lines
				// also set new contents to dst
				for cursor := 0; cursor < window; cursor++ {
					deletedLineIndex := deletedLinesIndices[cursor]
					addedLineIndex := addedLinesIndices[cursor]

					lines[deletedLineIndex].Operation = LineOperationModify
					lines[deletedLineIndex].Dst = lines[addedLineIndex].Dst
				}

				// remove lines with add operation that participated in modify
				if window > 0 && len(addedLinesIndices) != 0 {
					index := addedLinesIndices[0]
					lines = append(lines[:index], lines[index+window:]...)
				}

				// reset
				addedLinesIndices = nil
				deletedLinesIndices = nil
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
