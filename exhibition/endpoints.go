package exhibition

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type createPieceRequest struct {
	Piece Piece `json:"piece"`
}

type createPieceResponse struct {
	Created Piece  `json:"created"`
	Err     string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type getPiecesRequest struct {
}

type getPiecesResponse struct {
	Pieces []Piece `json:"data"`
	Err    string  `json:"err,omitempty"`
}

func makeCreatePieceEndpoint(svc ExhibitionService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createPieceRequest)
		v, err := svc.CreatePiece(req.Piece)
		if err != nil {
			return createPieceResponse{v, err.Error()}, nil
		}
		return createPieceResponse{v, ""}, nil
	}
}

func makeGetPiecesEndpoint(svc ExhibitionService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		// req := request.(getPiecesRequest)
		pieces, err := svc.GetPieces()
		if err != nil {
			return getPiecesResponse{pieces, err.Error()}, nil
		}
		return getPiecesResponse{pieces, ""}, nil
	}
}
