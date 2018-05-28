package rpc

import (
	"context"
	"github.com/coldze/primitives/custom_error"
	"github.com/coldze/primitives/json_rpc"
	"net/http"
	"github.com/coldze/memzie/interfaces/rpc/folders"
	"github.com/coldze/memzie/interfaces/rpc/words"
)

type HttpRequestHandler func(w http.ResponseWriter, r *http.Request)
type RegisteringFunc func(path string, handler HttpRequestHandler) custom_error.CustomError

type CreateParams struct {
	Field string `json:"field"`
}

type Echo struct {
	Received interface{} `json:"echo,omitempty"`
}

func Register(register RegisteringFunc) custom_error.CustomError {
	methods := map[string]json_rpc.HandlingInfo{
		"echo": json_rpc.HandlingInfo{
			Handle: func(ctx context.Context, request *json_rpc.RequestInfo) (json_rpc.ResponseInfo, json_rpc.ServerError) {
				v, ok := request.Data.(*CreateParams)
				if !ok {
					return json_rpc.ResponseInfo{
						Headers: nil,
						Data:    "FAILED",
					}, nil
				}
				return json_rpc.ResponseInfo{
					Headers: nil,
					Data: Echo{
						Received: v,
					},
				}, nil
			},
			NewParams: func() interface{} {
				return &CreateParams{}
			},
		},
	}

	err := register("system", json_rpc.CreateJSONRpcHandler(methods))
	if err != nil {
		return custom_error.NewErrorf(err, "Failed to register system rpc-methods.")
	}
	err = register("folders", json_rpc.CreateJSONRpcHandler(folders.NewMethods()))
	if err != nil {
		return custom_error.NewErrorf(err, "Failed to register folders rpc-methods.")
	}
	err = register("words", json_rpc.CreateJSONRpcHandler(words.NewMethods()))
	if err != nil {
		return custom_error.NewErrorf(err, "Failed to register words rpc-methods.")
	}
	return nil
}
