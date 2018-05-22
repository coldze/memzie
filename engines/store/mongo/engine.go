package mongo

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/mongo-go-driver/bson/objectid"
	"github.com/coldze/primitives/custom_error"
)

const boards_collection = "boards"

type CollectionFactory func(collectionName string) (store.Collection, custom_error.CustomError)

type BoardsFactory func(id string) (store.Boards, custom_error.CustomError)
type BoardFactory func(board *structs.Board) (store.Board, custom_error.CustomError)
type WordFactory func(word *structs.Word) (store.Word, custom_error.CustomError)

func NewBoardsFactory(factory BoardFactory, collFactory CollectionFactory) BoardsFactory {
	return func(id string) (store.Boards, custom_error.CustomError) {
		hexID, err := objectid.FromHex(id)
		if err != nil {
			return nil, custom_error.MakeErrorf("Failed to get boards. Client ID failed to convert. Error: %v", err)
		}
		coll, customErr := collFactory(boards_collection)
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to get '%v' collection.", boards_collection)
		}
		return &boardsImpl{
			boards:      coll,
			clientID:    hexID,
			factory:     factory,
		}, nil
	}
}

type root struct {
	collFactory CollectionFactory
	boards      store.Collection
	factory     BoardsFactory
}

func (b *root) GetBoards(clientID string) (store.Boards, custom_error.CustomError) {
	boards, err := b.factory(clientID)
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to get boards")
	}
	return boards, nil
}

type boardsImpl struct {
	clientID    objectid.ObjectID
	boards      store.Collection
	factory     BoardFactory
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

func NewRoot(collFactory CollectionFactory, factory BoardsFactory) (store.Root, custom_error.CustomError) {
	c, err := collFactory(boards_collection)
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to create new root.")
	}
	return &root{
		collFactory: collFactory,
		boards:      c,
		factory:     factory,
	}, nil
}
