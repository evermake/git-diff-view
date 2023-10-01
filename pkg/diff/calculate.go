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
	"github.com/samber/lo"
)

func Calculate(
	ctx context.Context,
	repoPath string,
	commitA, commitB string,
) ([]*Diff, error) {
	exists, err := gitutil.RevisionExists(ctx, repoPath, commitA)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("revision %s does not exists", commitA)
	}

	exists, err = gitutil.RevisionExists(ctx, repoPath, commitB)
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

		var lines []*Line
		for _, fragment := range file.TextFragments {
			var (
				srcLineNumber = fragment.OldPosition
				dstLineNumber = fragment.NewPosition
			)

			for _, fragmentLine := range fragment.Lines {
				line := new(Line)

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
					// as per request from Vlad =)
					//line.Dst.Content = fragmentLine.Line
					line.Dst.Number = dstLineNumber

					line.Src.Number = srcLineNumber
					line.Src.Content = fragmentLine.Line

					dstLineNumber++
					srcLineNumber++
				}

				lines = append(lines, line)
			}
		}

		var (
			deletedLines []*Line
			addedLines   []*Line
		)

		assignModifyOperations := func() {
			window := min(len(deletedLines), len(addedLines))

			// update operation to modify for deleted lines
			// also set new contents to dst
			for cursor := 0; cursor < window; cursor++ {
				deletedLine := deletedLines[cursor]
				addedLine := addedLines[cursor]

				deletedLine.Operation = LineOperationModify
				deletedLine.Dst = addedLine.Dst
			}

			// remove lines with add operation that participated in modify
			if window > 0 && len(addedLines) != 0 {
				_, index, _ := lo.FindIndexOf(lines, func(line *Line) bool {
					return line == addedLines[0]
				})

				if index+window >= len(lines) {
					lines = lines[:index]
				} else {
					lines = append(lines[:index], lines[index+window:]...)
				}
			}
		}

		for _, line := range lines {
			switch line.Operation {
			case LineOperationAdd:
				addedLines = append(addedLines, line)
			case LineOperationDelete:
				deletedLines = append(deletedLines, line)
			default:
				assignModifyOperations()

				// reset
				addedLines = nil
				deletedLines = nil
			}
		}
		assignModifyOperations()

		diff := diffs[name]
		diff.Lines = lines

		if diff.Status.Type == StatusAdd {
			diff.IsBinary, err = gitutil.IsBinary(ctx, repoPath, commitB, diff.Dst.Path)
			if err != nil {
				return nil, err
			}
		} else {
			diff.IsBinary, err = gitutil.IsBinary(ctx, repoPath, commitA, diff.Src.Path)
			if err != nil {
				return nil, err
			}
		}

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
