package order

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func MakeAddOrderEndpoint(o OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(OrderRequest)
		if err := o.Add(req); err != nil {
			return nil,err
		}
		return OrderResponse{"0","ok"},nil
	}
}

func DecodeOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

