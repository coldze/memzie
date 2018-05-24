package mongo2

import (
	"github.com/coldze/memzie/engines/logic"
	"github.com/coldze/memzie/engines/store"
	mgo_store "github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/impls"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/mongo-go-driver/bson"
	"github.com/coldze/mongo-go-driver/bson/objectid"
	"github.com/coldze/mongo-go-driver/core/options"
	mgo "github.com/coldze/mongo-go-driver/mongo"
	"github.com/coldze/primitives/custom_error"
)

type logicCollection struct {
	underlying *mgo_store.CollectionWrap
	sort       *bson.Document
	clientID   objectid.ObjectID
	boardID    objectid.ObjectID
	decoder    store.Decoder
	factory    mgo_store.WordFactory
}

func (c *logicCollection) Get(id string) (store.Word, custom_error.CustomError) {
	unhexWordID, err := objectid.FromHex(id)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert word ID. Error: %v", err)
	}
	filter := map[string]interface{}{
		"client_id": c.clientID,
		"board_id":  c.boardID,
		"_id":       unhexWordID,
	}
	res, customErr := c.underlying.FindOne(c.decoder, filter)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to FindOne.")
	}
	typed, ok := res.(store.Word)
	if !ok {
		return nil, custom_error.MakeErrorf("Failed to convert to type.")
	}
	return typed, nil
}

func (c *logicCollection) Next() (store.Word, custom_error.CustomError) {
	filter := map[string]interface{}{
		"client_id": c.clientID,
		"board_id":  c.boardID,
	}
	res, customErr := c.underlying.FindOne(c.decoder, filter, options.OptSort{
		Sort: c.sort,
	})
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to FindOne.")
	}
	typed, ok := res.(*structs.Word)
	if !ok {
		return nil, custom_error.MakeErrorf("Failed to convert to type.")
	}
	word, customErr := c.factory(typed)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to create word wrap.")
	}
	return word, nil
}

func NewLogic(client *mgo.Client, dbName string, collectionName string, boardID string, clientID string, factory mgo_store.WordFactory) (logic.Logic, custom_error.CustomError) {
	wrap, customErr := mgo_store.NewCollectionWrap(client, dbName, collectionName)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to create wrap.")
	}
	unhexBoardID, err := objectid.FromHex(boardID)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert board ID. Error: %v", err)
	}
	unhexClientID, err := objectid.FromHex(clientID)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert client ID. Error: %v", err)
	}
	return &logicCollection{
		clientID:   unhexClientID,
		boardID:    unhexBoardID,
		underlying: wrap,
		sort:       bson.NewDocument(bson.EC.Int64("weight", -1)),
		decoder:    impls.DecodeWord,
		factory:    factory,
	}, nil
}
