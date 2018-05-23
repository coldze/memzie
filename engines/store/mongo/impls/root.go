package impls

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/primitives/custom_error"
)

type root struct {
	collFactory mongo.CollectionFactory
	boards      store.Collection
	factory     mongo.BoardsFactory
}

func (b *root) GetBoards(clientID string) (store.Boards, custom_error.CustomError) {
	boards, err := b.factory(clientID)
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to get boards")
	}
	return boards, nil
}

func NewRoot(collFactory mongo.CollectionFactory, factory mongo.BoardsFactory) (store.Root, custom_error.CustomError) {
	c, err := collFactory(BOARDS_COLLECTION)
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to create new root.")
	}
	return &root{
		collFactory: collFactory,
		boards:      c,
		factory:     factory,
	}, nil
}
