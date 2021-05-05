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
		decodeGetPiecesRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Handle("/exhibition/v1/pieces", getPiecesHandler).Methods("GET")
	r.Handle("/exhibition/v1/pieces/{id}", createPieceHandler).Methods("POST")
	return r
}

func decodeCreatePieceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createPieceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGetPiecesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getPiecesRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
