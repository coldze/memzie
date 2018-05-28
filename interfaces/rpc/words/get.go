package words

import (
	"errors"
	"context"

	"github.com/coldze/primitives/json_rpc"
)

func decodeGetParams() interface{} {
	return nil
}

func NewGetHandler() json_rpc.RequestHandler {
	return func(ctx context.Context, request *json_rpc.RequestInfo) (json_rpc.ResponseInfo, json_rpc.ServerError) {
		return json_rpc.ResponseInfo{}, json_rpc.MakeError(0, 0, "NOT IMPLEMENTED", errors.New("Not implemented"))
	}
}
