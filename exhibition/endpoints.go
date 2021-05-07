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

type createArtistRequest struct {
	Artist Artist `json:"artist"`
}

type createArtistResponse struct {
	Created Artist `json:"created"`
	Err     string `json:"err,omitempty"`
}

type createSpaceRequest struct {
	Space Space `json:"space"`
}

type createSpaceResponse struct {
	Created Space  `json:"created"`
	Err     string `json:"err,omitempty"`
}

type getSpacesResponse struct {
	Spaces []Space `json:"data"`
	Err    string  `json:"err,omitempty"`
}

type getPiecesResponse struct {
	Pieces []Piece `json:"data"`
	Err    string  `json:"err,omitempty"`
}

type getArtistsResponse struct {
	Artists []Artist `json:"data"`
	Err     string   `json:"err,omitempty"`
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

func makeCreateArtistEndpoint(svc ExhibitionService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createArtistRequest)
		v, err := svc.CreateArtist(req.Artist)
		if err != nil {
			return createArtistResponse{v, err.Error()}, nil
		}
		return createArtistResponse{v, ""}, nil
	}
}

func makeCreateSpaceEndpoint(svc ExhibitionService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createSpaceRequest)
		v, err := svc.CreateSpace(req.Space)
		if err != nil {
			return createSpaceResponse{v, err.Error()}, nil
		}
		return createSpaceResponse{v, ""}, nil
	}
}

func makeGetPiecesEndpoint(svc ExhibitionService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		pieces, err := svc.GetPieces()
		if err != nil {
			return getPiecesResponse{pieces, err.Error()}, nil
		}
		return getPiecesResponse{pieces, ""}, nil
	}
}

func makeGetArtistsEndpoint(svc ExhibitionService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		artists, err := svc.GetArtists()
		if err != nil {
			return getArtistsResponse{artists, err.Error()}, nil
		}
		return getArtistsResponse{artists, ""}, nil
	}
}

func makeGetSpacesEndpoint(svc ExhibitionService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		spaces, err := svc.GetSpaces()
		if err != nil {
			return getSpacesResponse{spaces, err.Error()}, nil
		}
		return getSpacesResponse{spaces, ""}, nil
	}
}
