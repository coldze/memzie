package impls

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/primitives/custom_error"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type folderImpl struct {
	data    *structs.Folder
	words   store.Collection
	factory mongo.WordFactory
}

func DecodeWord(decode func(interface{}) error) (interface{}, custom_error.CustomError) {
	res := structs.Word{}
	err := decode(&res)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to decode word struct. Error: %v", err)
	}
	return &res, nil
}

func (b *folderImpl) GetName() string {
	return b.data.Name
}

func (b *folderImpl) GetDescription() string {
	return b.data.Description
}
func (b *folderImpl) GetID() string {
	return b.data.ID.Hex()
}

func (b *folderImpl) List(handle store.WordListHandle) custom_error.CustomError {
	collectFolders := func(decoded interface{}) (bool, custom_error.CustomError) {
		typed, ok := decoded.(*structs.Word)
		if !ok {
			return false, custom_error.MakeErrorf("Failed to convert input. Unexpected type: %T", decoded)
		}
		obj, err := b.factory(typed)
		if err != nil {
			return false, custom_error.NewErrorf(err, "Failed to create word wrap.")
		}
		next, err := handle(obj)
		if err != nil {
			return false, custom_error.NewErrorf(err, "Failed to handle word wrap.")
		}
		return next, nil
	}
	err := b.words.FindAll(DecodeWord, map[string]interface{}{"client_id": b.data.ClientID}, collectFolders)
	if err == nil {
		return nil
	}
	return custom_error.NewErrorf(err, "Failed to list words.")
}

func (b *folderImpl) Get(id string) (store.Word, custom_error.CustomError) {
	unhexID, err := objectid.FromHex(id)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert id. Error: %v", err)
	}
	filter := map[string]interface{}{
		"_id":       unhexID,
		"client_id": b.data.ClientID,
	}
	decoded, customErr := b.words.FindOne(DecodeWord, filter)
	if err != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to get word from collection.")
	}
	typed, ok := decoded.(*structs.Word)
	if !ok {
		return nil, custom_error.MakeErrorf("Failed to decode. Unexpected type: %T", decoded)
	}
	res, customErr := b.factory(typed)
	if customErr == nil {
		return res, nil
	}
	return nil, custom_error.NewErrorf(customErr, "Failed to create wrap.")
}

func (b *folderImpl) Create(params *store.WordCreateParams) (store.Word, custom_error.CustomError) {
	if params == nil {
		return nil, custom_error.MakeErrorf("Failed to create word. Empty params.")
	}
	object := structs.Word{
		ClientID:     b.data.ClientID,
		FolderID:     b.data.ID,
		Text:         params.Text,
		Translations: params.Translations,
		ID:           objectid.New(),
		Weight:       1.0,
	}
	id, customErr := b.words.Create(object)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to create word '%v'", *params)
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
	return nil, custom_error.NewErrorf(customErr, "Failed to create wrap word's object.")
}

func (b *folderImpl) Delete(id string) custom_error.CustomError {
	unhexID, err := objectid.FromHex(id)
	if err != nil {
		return custom_error.MakeErrorf("Failed to convert id. Error: %v", err)
	}
	filter := map[string]interface{}{
		"_id":       unhexID,
		"client_id": b.data.ClientID,
	}
	_, customErr := b.words.Delete(filter)
	if customErr == nil {
		return nil
	}
	return custom_error.NewErrorf(customErr, "Failed to delete word: %v", id)
}

func NewFolderFactory(factory mongo.WordFactory, collFactory mongo.CollectionFactory) mongo.FolderFactory {
	return func(folder *structs.Folder) (store.Folder, custom_error.CustomError) {
		if folder == nil {
			return nil, custom_error.MakeErrorf("Failed to wrap folder. Data is nil.")
		}
		coll, customErr := collFactory(WORDS_COLLECTION)
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to get '%v' collection.", WORDS_COLLECTION)
		}
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to create folder wrap.")
		}
		return &folderImpl{
			data:    folder,
			words:   coll,
			factory: factory,
		}, nil
	}
}
