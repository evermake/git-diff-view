package gitutil

import (
	"context"
	"fmt"
	"os/exec"
)

func ReadFile(ctx context.Context, revision, path string) ([]byte, error) {
	exists, err := RevisionExists(ctx, revision)
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

	return cmd.Output()
}

func RevisionExists(ctx context.Context, revision string) (bool, error) {
	cmd := exec.CommandContext(
		ctx,
		"git",
		"cat-file",
		"-e",
		fmt.Sprintf("%s^{commit}", revision),
	)

	err := cmd.Run()

	if err == nil {
		return true, nil
	}

	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 128 {
		return false, nil
	}

	return false, err
}
