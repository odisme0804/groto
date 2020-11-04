package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"groto/internal/endpoint"
	"groto/pkg"
)

func NewHTTPServer(e endpoint.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(errorEncoder),
	}

	r.Methods("Get").Path("/hello").Handler(kithttp.NewServer(
		e.SayHello,
		decodeSayHelloRequest,
		encodeSayHelloResponse,
		options...,
	))

	return r
}

func decodeSayHelloRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req pkg.HelloRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}

	return &req, nil
}

func encodeSayHelloResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(429)
	w.Write([]byte(err.Error()))
}
