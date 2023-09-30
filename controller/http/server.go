package http

import (
	"context"

	"github.com/evermake/git-diff-view/controller/http/openapi"
)

var _ openapi.StrictServerInterface = (*Server)(nil)

type Server struct{}

func (s *Server) GetDiffMap(ctx context.Context, request openapi.GetDiffMapRequestObject) (openapi.GetDiffMapResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetDiffPart(ctx context.Context, request openapi.GetDiffPartRequestObject) (openapi.GetDiffPartResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
