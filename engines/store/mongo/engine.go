package mongo

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/primitives/custom_error"
)

type CollectionFactory func(collectionName string) (store.Collection, custom_error.CustomError)

type FoldersFactory func(id string) (store.Folders, custom_error.CustomError)
type FolderFactory func(folder *structs.Folder) (store.Folder, custom_error.CustomError)
type WordFactory func(word *structs.Word) (store.Word, custom_error.CustomError)
