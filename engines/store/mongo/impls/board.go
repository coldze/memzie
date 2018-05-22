package impls

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/mongo-go-driver/bson/objectid"
	"github.com/coldze/primitives/custom_error"
)

const words_collection = "words"

type boardImpl struct {
	data        *structs.Board
	words       store.Collection
	factory     mongo.WordFactory
}

func decodeWord(decode func(interface{}) error) (interface{}, custom_error.CustomError) {
	res := structs.Word{}
	err := decode(&res)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to decode word struct. Error: %v", err)
	}
	return &res, nil
}

func (b *boardImpl) GetName() string {
	return b.data.Name
}

func (b *boardImpl) GetDescription() string {
	return b.data.Description
}
func (b *boardImpl) GetID() string {
	return b.data.ID.Hex()
}

func (b *boardImpl) List(handle store.WordListHandle) custom_error.CustomError {
	collectBoards := func(decoded interface{}) (bool, custom_error.CustomError) {
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
	err := b.words.FindAll(decodeWord, map[string]interface{}{"client_id": b.data.ClientID}, collectBoards)
	if err == nil {
		return nil
	}
	return custom_error.NewErrorf(err, "Failed to list words.")
}

func (b *boardImpl) Get(id string) (store.Word, custom_error.CustomError) {
	unhexID, err := objectid.FromHex(id)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert id. Error: %v", err)
	}
	filter := map[string]interface{}{
		"_id":       unhexID,
		"client_id": b.data.ClientID,
	}
	decoded, customErr := b.words.FindOne(decodeWord, filter)
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

func (b *boardImpl) Create(params *store.WordCreateParams) (store.Word, custom_error.CustomError) {
	if params == nil {
		return nil, custom_error.MakeErrorf("Failed to create word. Empty params.")
	}
	object := structs.Word{
		ClientID:     b.data.ClientID,
		BoardID:      b.data.ID,
		Text:         params.Text,
		Translations: params.Translations,
		ID:           objectid.New(),
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

func (b *boardImpl) Delete(id string) custom_error.CustomError {
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

func NewBoardFactory(factory mongo.WordFactory, collFactory mongo.CollectionFactory) mongo.BoardFactory {
	return func(board *structs.Board) (store.Board, custom_error.CustomError) {
		if board == nil {
			return nil, custom_error.MakeErrorf("Failed to wrap board. Data is nil.")
		}
		coll, customErr := collFactory(words_collection)
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to get '%v' collection.", words_collection)
		}
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to create board wrap.")
		}
		return &boardImpl{
			data:    board,
			words:   coll,
			factory: factory,
		}, nil
	}
}
