package gitutil

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func Branches(ctx context.Context, repo string) ([]Branch, error) {
	cmd := exec.CommandContext(
		ctx,
		"git",
		"branch",
		"--format",
		"%(refname:short)",
	)
	cmd.Dir = repo

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	names := bytes.Fields(stdout)

	branches := make([]Branch, len(names))
	for i, name := range names {
		branches[i] = Branch{Name: string(name)}
	}

	return branches, nil
}

func BranchCommits(ctx context.Context, repo string, branch string) ([]CommitPreview, error) {
	cmd := exec.CommandContext(
		ctx,
		"git",
		"log",
		"--format=oneline",
		branch,
	)

	cmd.Dir = repo

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	stdout = bytes.TrimSpace(stdout)
	lines := bytes.Split(stdout, []byte("\n"))

	commits := make([]CommitPreview, len(lines))
	for i, line := range lines {
		fields := bytes.Fields(line)

		commits[i] = CommitPreview{
			SHA1:    fields[0],
			Message: string(bytes.Join(fields[1:], []byte(" "))),
		}
	}

	return commits, nil
}

func ReadFile(ctx context.Context, repo string, revision, path string) ([]byte, error) {
	exists, err := RevisionExists(ctx, repo, revision)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("revision %s does not exist", revision)
	}

	cmd := exec.CommandContext(
		ctx,
		"git",
		"show",
		fmt.Sprintf("%s:%s", revision, path),
	)

	cmd.Dir = repo

	return cmd.Output()
}

func RevisionExists(ctx context.Context, repo string, revision string) (bool, error) {
	cmd := exec.CommandContext(
		ctx,
		"git",
		"cat-file",
		"-e",
		fmt.Sprintf("%s^{commit}", revision),
	)
	cmd.Dir = repo

	err := cmd.Run()

	if err == nil {
		return true, nil
	}

	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 128 {
		return false, nil
	}

	return false, err
}
