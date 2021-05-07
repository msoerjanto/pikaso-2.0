package exhibition

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandlers(es ExhibitionService, logger kitlog.Logger) http.Handler {
	createPieceHandler := kithttp.NewServer(
		makeCreatePieceEndpoint(es),
		decodeCreatePieceRequest,
		encodeResponse,
	)

	getPiecesHandler := kithttp.NewServer(
		makeGetPiecesEndpoint(es),
		ignoreRequestDecoder,
		encodeResponse,
	)

	createArtistHandler := kithttp.NewServer(
		makeCreateArtistEndpoint(es),
		decodeCreateArtistRequest,
		encodeResponse,
	)

	getArtistsHandler := kithttp.NewServer(
		makeGetArtistsEndpoint(es),
		ignoreRequestDecoder,
		encodeResponse,
	)

	createSpaceHandler := kithttp.NewServer(
		makeCreateSpaceEndpoint(es),
		decodeCreateSpaceRequest,
		encodeResponse,
	)

	getSpacesHandler := kithttp.NewServer(
		makeGetSpacesEndpoint(es),
		ignoreRequestDecoder,
		encodeResponse,
	)

	r := mux.NewRouter()

	r.Handle("/exhibition/v1/pieces", getPiecesHandler).Methods("GET")
	r.Handle("/exhibition/v1/pieces", createPieceHandler).Methods("POST")

	r.Handle("/exhibition/v1/artists", getArtistsHandler).Methods("GET")
	r.Handle("/exhibition/v1/artists", createArtistHandler).Methods("POST")

	r.Handle("/exhibition/v1/spaces", getSpacesHandler).Methods("GET")
	r.Handle("/exhibition/v1/spaces", createSpaceHandler).Methods("POST")

	return r
}

func decodeCreatePieceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createPieceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCreateArtistRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createArtistRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCreateSpaceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createSpaceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func ignoreRequestDecoder(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
