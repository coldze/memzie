package impls

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/primitives/custom_error"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type boardsImpl struct {
	clientID objectid.ObjectID
	boards   store.Collection
	factory  mongo.BoardFactory
}

func decodeBoard(decode func(interface{}) error) (interface{}, custom_error.CustomError) {
	res := structs.Board{}
	err := decode(&res)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to decode board struct. Error: %v", err)
	}
	return &res, nil
}

func (b *boardsImpl) List(handle store.BoardListHandle) custom_error.CustomError {
	collectBoards := func(decoded interface{}) (bool, custom_error.CustomError) {
		typed, ok := decoded.(*structs.Board)
		if !ok {
			return false, custom_error.MakeErrorf("Failed to convert input. Unexpected type: %T", decoded)
		}
		obj, err := b.factory(typed)
		if err != nil {
			return false, custom_error.NewErrorf(err, "Failed to create board wrap.")
		}
		next, err := handle(obj)
		if err != nil {
			return false, custom_error.NewErrorf(err, "Failed to handle board wrap.")
		}
		return next, nil
	}
	err := b.boards.FindAll(decodeBoard, map[string]interface{}{"client_id": b.clientID}, collectBoards)
	if err == nil {
		return nil
	}
	return custom_error.NewErrorf(err, "Failed to list boards.")
}

func (b *boardsImpl) Get(id string) (store.Board, custom_error.CustomError) {
	unhexID, err := objectid.FromHex(id)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert id. Error: %v", err)
	}
	filter := map[string]interface{}{
		"_id":       unhexID,
		"client_id": b.clientID,
	}
	decoded, customErr := b.boards.FindOne(decodeBoard, filter)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to get board from collection.")
	}
	typed, ok := decoded.(*structs.Board)
	if !ok {
		return nil, custom_error.MakeErrorf("Failed to decode. Unexpected type: %T", decoded)
	}
	res, customErr := b.factory(typed)
	if customErr == nil {
		return res, nil
	}
	return nil, custom_error.NewErrorf(customErr, "Failed to create wrap.")
}

func (b *boardsImpl) Create(name string, description string) (store.Board, custom_error.CustomError) {
	object := structs.Board{
		ClientID:    b.clientID,
		Name:        name,
		Description: description,
		ID:          objectid.New(),
	}
	id, customErr := b.boards.Create(object)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to create board '%v'", name)
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
	return nil, custom_error.NewErrorf(customErr, "Failed to create wrap board's object.")
}

func (b *boardsImpl) Delete(id string) custom_error.CustomError {
	unhexID, err := objectid.FromHex(id)
	if err != nil {
		return custom_error.MakeErrorf("Failed to convert id. Error: %v", err)
	}
	filter := map[string]interface{}{
		"_id":       unhexID,
		"client_id": b.clientID,
	}
	_, customErr := b.boards.Delete(filter)
	if customErr == nil {
		return nil
	}
	return custom_error.NewErrorf(customErr, "Failed to delete board: %v", id)
}

func NewBoardsFactory(factory mongo.BoardFactory, collFactory mongo.CollectionFactory) mongo.BoardsFactory {
	return func(id string) (store.Boards, custom_error.CustomError) {
		hexID, err := objectid.FromHex(id)
		if err != nil {
			return nil, custom_error.MakeErrorf("Failed to get boards. Client ID failed to convert. Error: %v", err)
		}
		coll, customErr := collFactory(BOARDS_COLLECTION)
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to get '%v' collection.", BOARDS_COLLECTION)
		}
		return &boardsImpl{
			boards:   coll,
			clientID: hexID,
			factory:  factory,
		}, nil
	}
}
