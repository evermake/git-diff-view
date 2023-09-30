package diff

import (
	"bytes"
	"context"
	"os/exec"
)

func Calculate(
	ctx context.Context,
	repoPath string,
	commitA, commitB string,
) ([]Diff, error) {
	cmd := exec.CommandContext(
		ctx,
		"git",
		"diff",
		"--patch-with-raw",
		commitA,
		commitB,
	)

	cmd.Dir = repoPath

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	stdout = bytes.TrimSpace(stdout)
	lines := bytes.Split(stdout, []byte("\n"))

	diffs := make(map[string]Diff)

	var (
		patchPath string
		patch     [][]byte
	)

	for _, line := range lines {
		if bytes.HasPrefix(line, []byte(":")) {
			diff, err := parseDiff(line)
			if err != nil {
				return nil, err
			}

			diffs[diff.Src.Path] = diff
		} else {
			if bytes.HasPrefix(line, []byte("diff --git")) {
				if patchPath != "" {
					diff := diffs[patchPath]
					diff.Patch = bytes.Join(patch, []byte("\n"))
					diffs[patchPath] = diff
				}

				fields := bytes.Fields(line)
				if len(fields) < 3 {
					return nil, ErrMalformedDiff
				}

				patchPath = string(bytes.TrimPrefix(fields[2], []byte("a/")))
				patch = nil
			}

			patch = append(patch, line)
		}
	}

	if patchPath != "" {
		diff := diffs[patchPath]
		diff.Patch = bytes.Join(patch, []byte("\n"))
		diffs[patchPath] = diff
	}

	diffsSlice := make([]Diff, 0, len(diffs))
	for _, diff := range diffs {
		diffsSlice = append(diffsSlice, diff)
	}

	return diffsSlice, nil
}
