package v1

import (
	"context"
	"math"
	"strings"

	"github.com/evermake/git-diff-view/internal/controller/http/v1/openapi"
	"github.com/evermake/git-diff-view/pkg/diff"
	"github.com/evermake/git-diff-view/pkg/gitutil"
	"github.com/jellydator/ttlcache/v3"
	"github.com/samber/lo"
)

var _ openapi.StrictServerInterface = (*Server)(nil)

func NewServer(repoPath string) *Server {
	return &Server{
		repoPath: repoPath,
		diffCache: ttlcache.New[string, []*combinedDiff](
			ttlcache.WithCapacity[string, []*combinedDiff](10),
		),
		fileCache: ttlcache.New[string, []string](
			ttlcache.WithCapacity[string, []string](10),
		),
	}
}

type Server struct {
	repoPath  string
	fileCache *ttlcache.Cache[string, []string]
	diffCache *ttlcache.Cache[string, []*combinedDiff]
}

func (s *Server) getDiffs(ctx context.Context, commitA, commitB string) ([]*combinedDiff, error) {
	if entry := s.diffCache.Get(commitA + commitB); entry != nil {
		return entry.Value(), nil
	}

	diffs, err := diff.Calculate(ctx, s.repoPath, commitA, commitB)
	if err != nil {
		return nil, err
	}

	var lineNumber int

	entries := make([]*combinedDiff, len(diffs))
	for i, d := range diffs {
		entry := &combinedDiff{}
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
			end := start + len(d.Lines) - 1

			entry.FileDiff.Lines = openapi.Range{
				Start: start,
				End:   end,
			}

			lineNumber = end
		}

		entries[i] = entry
	}

	s.diffCache.Set(commitA+commitB, entries, ttlcache.DefaultTTL)
	return entries, nil
}

func (s *Server) GetFile(ctx context.Context, request openapi.GetFileRequestObject) (openapi.GetFileResponseObject, error) {
	revision := "HEAD"
	if request.Params.Revision != nil {
		revision = *request.Params.Revision
	}

	var (
		start, end = 0, int(math.Inf(1))
	)

	if request.Params.Start != nil {
		start = *request.Params.Start
	}

	if request.Params.End != nil {
		end = *request.Params.End
	}

	if start > end {
		return openapi.GetFile400JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "start is greater than end",
			},
		}, nil
	}

	if file := s.fileCache.Get(revision + request.Params.Path); file != nil {
		lines := file.Value()

		if start < 0 {
			start = 0
		}
		if start >= len(lines) {
			start = len(lines) - 1
		}

		if end < 0 {
			end = 0
		}
		if end >= len(lines) {
			end = len(lines)
		}

		return openapi.GetFile200JSONResponse(lines[start:end]), nil
	}

	contents, err := gitutil.ReadFile(ctx, revision, request.Params.Path)
	if err != nil {
		return openapi.GetFile400JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: err.Error(),
			},
		}, nil
	}

	s.fileCache.Set(
		revision+request.Params.Path,
		strings.Split(string(contents), "\n"),
		ttlcache.DefaultTTL,
	)

	return s.GetFile(ctx, request)
}

func (s *Server) GetDiffMap(ctx context.Context, request openapi.GetDiffMapRequestObject) (openapi.GetDiffMapResponseObject, error) {
	diffs, err := s.getDiffs(ctx, request.Params.A, request.Params.B)
	if err != nil {
		return openapi.GetDiffMap400JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: err.Error(),
			},
		}, nil
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
		return openapi.GetDiffPart400JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: err.Error(),
			},
		}, nil
	}

	start, end := request.Params.Start-1, request.Params.End-1

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

	if start < 0 || start >= len(lines) {
		return openapi.GetDiffPart400JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "start is out of bounds",
			},
		}, nil
	}

	if end < 0 || end >= len(lines) {
		return openapi.GetDiffPart400JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "end is out of bounds",
			},
		}, nil
	}

	linesDiff := lo.Map(lines[start:end+1], func(line diff.Line, _ int) openapi.LineDiff {
		var operation openapi.LineDiffOperation
		switch line.Operation {
		case diff.LineOperationModify:
			operation = openapi.LineDiffOperationM
		case diff.LineOperationAdd:
			operation = openapi.LineDiffOperationA
		case diff.LineOperationDelete:
			operation = openapi.LineDiffOperationD
		}

		return openapi.LineDiff{
			Operation: operation,
			Dst: openapi.LineState{
				Content: line.Dst.Content,
				Number:  line.Dst.Number,
			},
			Src: openapi.LineState{
				Content: line.Src.Content,
				Number:  line.Src.Number,
			},
		}
	})

	return openapi.GetDiffPart200JSONResponse(linesDiff), nil
}
