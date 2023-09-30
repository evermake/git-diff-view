package v1

import (
	"github.com/evermake/git-diff-view/internal/controller/http/v1/openapi"
	"github.com/evermake/git-diff-view/pkg/diff"
)

type combinedDiff struct {
	FileDiff *openapi.FileDiff
	Diff     *diff.Diff
}
