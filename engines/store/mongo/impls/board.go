package impls

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/primitives/custom_error"
)

const words_collection = "words"

type boardImpl struct {
	data       *structs.Board
	collection store.Collection
	factory    mongo.WordFactory
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
	return custom_error.MakeErrorf("Not implemented")
}

func (b *boardImpl) Get(id string) (store.Word, custom_error.CustomError) {
	return nil, custom_error.MakeErrorf("Not implemented")
}

func (b *boardImpl) Create(params *store.WordCreateParams) (store.Word, custom_error.CustomError) {
	return nil, custom_error.MakeErrorf("Not implemented")
}

func (b *boardImpl) Delete(id string) custom_error.CustomError {
	return custom_error.MakeErrorf("Not implemented")
}

func NewBoardFactory(factory mongo.WordFactory) mongo.BoardFactory {
	return func(board *structs.Board, engine store.Engine) (store.Board, custom_error.CustomError) {
		if board == nil {
			return nil, custom_error.MakeErrorf("Failed to wrap board. Data is nil.")
		}
		coll, customErr := engine.GetCollection(words_collection)
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to get '%v' collection.", words_collection)
		}
		if customErr != nil {
			return nil, custom_error.NewErrorf(customErr, "Failed to create board wrap.")
		}
		return &boardImpl{
			data:       board,
			collection: coll,
			factory:    factory,
		}, nil
	}
}
