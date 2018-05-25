package impls

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/primitives/custom_error"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type foldersImpl struct {
	clientID objectid.ObjectID
	folders  store.Collection
	factory  mongo.FolderFactory
}

func decodeFolder(decode func(interface{}) error) (interface{}, custom_error.CustomError) {
	res := structs.Folder{}
	err := decode(&res)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to decode folder struct. Error: %v", err)
	}
	return &res, nil
}

func (b *foldersImpl) List(handle store.FolderListHandle) custom_error.CustomError {
	collectFolders := func(decoded interface{}) (bool, custom_error.CustomError) {
		typed, ok := decoded.(*structs.Folder)
		if !ok {
			return false, custom_error.MakeErrorf("Failed to convert input. Unexpected type: %T", decoded)
		}
		obj, err := b.factory(typed)
		if err != nil {
			return false, custom_error.NewErrorf(err, "Failed to create folder wrap.")
		}
		next, err := handle(obj)
		if err != nil {
			return false, custom_error.NewErrorf(err, "Failed to handle folder wrap.")
		}
		return next, nil
	}
	err := b.folders.FindAll(decodeFolder, map[string]interface{}{"client_id": b.clientID}, collectFolders)
	if err == nil {
		return nil
	}
	return custom_error.NewErrorf(err, "Failed to list folders.")
}

func (b *foldersImpl) Get(id string) (store.Folder, custom_error.CustomError) {
	unhexID, err := objectid.FromHex(id)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert id. Error: %v", err)
	}
	filter := map[string]interface{}{
		"_id":       unhexID,
		"client_id": b.clientID,
	}
	decoded, customErr := b.folders.FindOne(decodeFolder, filter)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to get folder from collection.")
	}
	typed, ok := decoded.(*structs.Folder)
	if !ok {
		return nil, custom_error.MakeErrorf("Failed to decode. Unexpected type: %T", decoded)
	}
	res, customErr := b.factory(typed)
	if customErr == nil {
		return res, nil
	}
	return nil, custom_error.NewErrorf(customErr, "Failed to create wrap.")
}

func (b *foldersImpl) Create(name string, description string) (store.Folder, custom_error.CustomError) {
	object := structs.Folder{
		ClientID:    b.clientID,
		Name:        name,
		Description: description,
		ID:          objectid.New(),
	}
	id, customErr := b.folders.Create(object)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to create folder '%v'", name)
	}
	var err error
	object.ID, err = objectid.FromHex(id)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert new object id. Error: %v", err)
	}
	res, customErr := b.factory(&object)
	if customErr == nil {
		return res, nil
	}
	return nil, custom_error.NewErrorf(customErr, "Failed to create wrap folder's object.")
}

func (b *foldersImpl) Delete(id string) custom_error.CustomError {
	unhexID, err := objectid.FromHex(id)
	if err != nil {
		return custom_error.MakeErrorf("Failed to convert id. Error: %v", err)
	}
	filter := map[string]interface{}{
		"_id":       unhexID,
		"client_id": b.clientID,
	}
	_, customErr := b.folders.Delete(filter)
	if customErr == nil {
		return nil
	}
	return custom_error.NewErrorf(customErr, "Failed to delete folder: %v", id)
}

func NewFoldersFactory(factory mongo.FolderFactory, collFactory mongo.CollectionFactory) mongo.FoldersFactory {
	return func(id string) (store.Folders, custom_error.CustomError) {
		hexID, err := objectid.FromHex(id)
		if err != nil {
			return nil, custom_error.MakeErrorf("Failed to get folders. Client ID failed to convert. Error: %v", err)
		}
		coll, customErr := collFactory(FOLDERS_COLLECTION)
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to get '%v' collection.", FOLDERS_COLLECTION)
		}
		return &foldersImpl{
			folders:  coll,
			clientID: hexID,
			factory:  factory,
		}, nil
	}
}
