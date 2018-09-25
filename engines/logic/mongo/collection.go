package mongo2

import (
	"github.com/coldze/memzie/engines/logic"
	"github.com/coldze/memzie/engines/store"
	mgo_store "github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/impls"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/primitives/custom_error"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

type logicCollection struct {
	underlying *mgo_store.CollectionWrap
	sort       *bson.Document
	clientID   objectid.ObjectID
	folderID   objectid.ObjectID
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
		"folder_id": c.folderID,
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
		"folder_id": c.folderID,
	}
	opts := &findopt.OneBundle{}
	opts = opts.Sort(c.sort)
	res, customErr := c.underlying.FindOne(c.decoder, filter, opts)
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

func NewLogic(client *mgo.Client, dbName string, collectionName string, folderID string, clientID string, factory mgo_store.WordFactory) (logic.Logic, custom_error.CustomError) {
	wrap, customErr := mgo_store.NewCollectionWrap(client, dbName, collectionName)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to create wrap.")
	}
	unhexFolderID, err := objectid.FromHex(folderID)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert board ID. Error: %v", err)
	}
	unhexClientID, err := objectid.FromHex(clientID)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert client ID. Error: %v", err)
	}
	return &logicCollection{
		clientID:   unhexClientID,
		folderID:   unhexFolderID,
		underlying: wrap,
		sort:       bson.NewDocument(bson.EC.Int64("weight", -1)),
		decoder:    impls.DecodeWord,
		factory:    factory,
	}, nil
}
