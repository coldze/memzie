package logic

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/primitives/custom_error"
)

type Logic interface {
	Get(id string) (store.Word, custom_error.CustomError)
	Next() (store.Word, custom_error.CustomError)
}
