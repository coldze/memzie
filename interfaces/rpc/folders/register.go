package folders

import (
	"github.com/coldze/primitives/json_rpc"
)

func NewMethods() map[string]json_rpc.HandlingInfo {
	methods := map[string]json_rpc.HandlingInfo{
		"create": json_rpc.HandlingInfo{
			NewParams: decodeCreateParams,
			Handle:    NewCreateHandler(),
		},
		"delete": json_rpc.HandlingInfo{
			NewParams: decodeDeleteParams,
			Handle:    NewDeleteHandler(),
		},
		"get": json_rpc.HandlingInfo{
			NewParams: decodeGetParams,
			Handle:    NewGetHandler(),
		},
		"list": json_rpc.HandlingInfo{
			NewParams: decodeListParams,
			Handle:    NewListHandler(),
		},
	}
	return methods
}

