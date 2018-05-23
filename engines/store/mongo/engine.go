package mongo

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/primitives/custom_error"
)

type CollectionFactory func(collectionName string) (store.Collection, custom_error.CustomError)

type BoardsFactory func(id string) (store.Boards, custom_error.CustomError)
type BoardFactory func(board *structs.Board) (store.Board, custom_error.CustomError)
type WordFactory func(word *structs.Word) (store.Word, custom_error.CustomError)
