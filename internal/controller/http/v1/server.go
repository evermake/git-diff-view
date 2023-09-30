package v1

import (
	"context"
	"os"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/evermake/git-diff-view/internal/controller/http/v1/openapi"
	"github.com/evermake/git-diff-view/pkg/diff"
	"github.com/samber/lo"
)

var _ openapi.StrictServerInterface = (*Server)(nil)

func NewServer() *Server {
	return &Server{
		cache: newDiffCache(),
	}
}

type Server struct {
	cache *diffCache
}

func (s *Server) getDiffs(ctx context.Context, commitA, commitB string) ([]*diffCacheEntry, error) {
	if entry, ok := s.cache.Get(commitA, commitB); ok {
		return entry, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	diffs, err := diff.Calculate(ctx, wd, commitA, commitB)
	if err != nil {
		return nil, err
	}

	var lineNumber int

	entries := make([]*diffCacheEntry, len(diffs))
	for i, d := range diffs {
		entry := &diffCacheEntry{}
		entry.Diff = d
		entry.FileDiff = &openapi.FileDiff{
			IsBinary: d.IsBinary,
			Src: openapi.State{
				Path: d.Src.Path,
			},
			Dst: openapi.State{
				Path: d.Dst.Path,
			},
			Status: openapi.Status{
				Score: d.Status.Score,
				Type:  openapi.StatusType(d.Status.Type),
			},
		}

		{
			start := lineNumber + 1
			end := start + len(d.Lines)

			entry.FileDiff.Lines = openapi.Range{
				Start: start,
				End:   end,
			}

			lineNumber = end
		}

		entries[i] = entry
	}

	s.cache.Set(commitA, commitB, entries)
	return entries, nil
}

func (s *Server) GetDiffMap(ctx context.Context, request openapi.GetDiffMapRequestObject) (openapi.GetDiffMapResponseObject, error) {
	diffs, err := s.getDiffs(ctx, request.Params.A, request.Params.B)
	if err != nil {
		return nil, err
	}

	files := make([]openapi.FileDiff, len(diffs))
	var lastLine int
	for i, d := range diffs {
		files[i] = *d.FileDiff
		lastLine = d.FileDiff.Lines.End
	}

	return openapi.GetDiffMap200JSONResponse{
		Files:      files,
		LinesTotal: lastLine,
	}, nil
}

func (s *Server) GetDiffPart(ctx context.Context, request openapi.GetDiffPartRequestObject) (openapi.GetDiffPartResponseObject, error) {
	diffs, err := s.getDiffs(ctx, request.Params.A, request.Params.B)
	if err != nil {
		return nil, err
	}

	start, end := request.Params.Start, request.Params.End

	if start > end {
		return openapi.GetDiffPart400JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "start is greater than end",
			},
		}, nil
	}

	var lines []diff.Line
	for _, d := range diffs {
		lines = append(lines, d.Diff.Lines...)
	}

	if start <= 0 || start >= len(lines) {
		return openapi.GetDiffPart400JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "start is out of bounds",
			},
		}, nil
	}

	if end <= 0 || end >= len(lines) {
		return openapi.GetDiffPart400JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "end is out of bounds",
			},
		}, nil
	}

	linesDiff := lo.Map(lines[start:end], func(line diff.Line, _ int) openapi.LineDiff {
		var op openapi.LineDiffOp
		switch line.Op {
		case gitdiff.OpAdd:
			op = openapi.LineDiffOpA
		case gitdiff.OpDelete:
			op = openapi.LineDiffOpD
		}

		return openapi.LineDiff{
			Content:       line.Content,
			SrcLineNumber: line.NumberInSrc,
			DstLineNumber: line.NumberInDst,
			Op:            op,
		}
	})

	return openapi.GetDiffPart200JSONResponse(linesDiff), nil
}
