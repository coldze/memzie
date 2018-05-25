package impls

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/primitives/custom_error"
)

type root struct {
	collFactory mongo.CollectionFactory
	folders     store.Collection
	factory     mongo.FoldersFactory
}

func (b *root) GetFolders(clientID string) (store.Folders, custom_error.CustomError) {
	folders, err := b.factory(clientID)
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to get folders")
	}
	return folders, nil
}

func NewRoot(collFactory mongo.CollectionFactory, factory mongo.FoldersFactory) (store.Root, custom_error.CustomError) {
	c, err := collFactory(FOLDERS_COLLECTION)
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to create new root.")
	}
	return &root{
		collFactory: collFactory,
		folders:     c,
		factory:     factory,
	}, nil
}
