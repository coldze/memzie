package mongo

import (
	"context"

	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/mongo-go-driver/bson/objectid"
	mgo "github.com/coldze/mongo-go-driver/mongo"
	"github.com/coldze/primitives/custom_error"
)

const boards_collection = "boards"

type engine struct {
	client *mgo.Client
	dbName string
}

func (e *engine) GetCollection(name string) (store.Collection, custom_error.CustomError) {
	c, err := NewCollection(e.client, e.dbName, name)
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to get collection '%v'", name)
	}
	return c, nil
}

type BoardsFactory func(id string, engine store.Engine) (store.Boards, custom_error.CustomError)
type BoardFactory func(board *structs.Board, engine store.Engine) (store.Board, custom_error.CustomError)
type WordFactory func(word *structs.Word, engine store.Engine) (store.Word, custom_error.CustomError)

func NewBoardsFactory(factory BoardFactory) BoardsFactory {
	return func(id string, engine store.Engine) (store.Boards, custom_error.CustomError) {
		hexID, err := objectid.FromHex(id)
		if err != nil {
			return nil, custom_error.MakeErrorf("Failed to get boards. Client ID failed to convert. Error: %v", err)
		}
		coll, customErr := engine.GetCollection(boards_collection)
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to get '%v' collection.", boards_collection)
		}
		return &boardsImpl{
			engine:   engine,
			boards:   coll,
			clientID: hexID,
			factory:  factory,
		}, nil
	}
}

type root struct {
	engine  store.Engine
	boards  store.Collection
	factory BoardsFactory
}

func (b *root) GetBoards(clientID string) (store.Boards, custom_error.CustomError) {
	boards, err := b.factory(clientID, b.engine)
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to get boards")
	}
	return boards, nil
}

func (e *engine) Close() {
	e.client.Disconnect(context.Background())
}

type boardsImpl struct {
	clientID objectid.ObjectID
	boards   store.Collection
	engine   store.Engine
	factory  BoardFactory
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
		obj, err := b.factory(typed, b.engine)
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
	decoded, err := b.boards.FindOne(decodeBoard, id, map[string]interface{}{"client_id": b.clientID})
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to get board from collection.")
	}
	typed, ok := decoded.(*structs.Board)
	if !ok {
		return nil, custom_error.MakeErrorf("Failed to decode. Unexpected type: %T", decoded)
	}
	res, customErr := b.factory(typed, b.engine)
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
	res, customErr := b.factory(&object, b.engine)
	if customErr == nil {
		return res, nil
	}
	return nil, custom_error.NewErrorf(customErr, "Failed to create wrap board's object.")
}

func (b *boardsImpl) Delete(id string) custom_error.CustomError {
	_, err := b.boards.Delete(id, map[string]interface{}{"client_id": b.clientID})
	if err == nil {
		return nil
	}
	return custom_error.NewErrorf(err, "Failed to delete board: %v", id)
}

func NewEngine(url string, dbName string) (store.Engine, custom_error.CustomError) {
	ctx := context.Background()
	client, err := mgo.Connect(ctx, url, nil)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to connect to mongo-db. Error: %v", err)
	}
	if client == nil {
		return nil, custom_error.MakeErrorf("Something went wrong. Client is nil")
	}
	return &engine{
		client: client,
		dbName: dbName,
	}, nil
}

func NewRoot(engine store.Engine, factory BoardsFactory) (store.Root, custom_error.CustomError) {
	c, err := engine.GetCollection(boards_collection)
	if err != nil {
		return nil, custom_error.NewErrorf(err, "Failed to create new root.")
	}
	return &root{
		engine:  engine,
		boards:  c,
		factory: factory,
	}, nil
}
