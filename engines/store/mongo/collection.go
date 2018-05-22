package mongo

import (
	"context"

	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/mongo-go-driver/bson"
	"github.com/coldze/mongo-go-driver/bson/objectid"
	mgo "github.com/coldze/mongo-go-driver/mongo"
	"github.com/coldze/primitives/custom_error"
)

type collection struct {
	dbName     string
	collName   string
	collection *mgo.Collection
}

func (c *collection) FindOne(decoder store.Decoder, id string, additionalFilters map[string]interface{}) (interface{}, custom_error.CustomError) {
	unhexID, err := objectid.FromHex(id)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert id. Error: %v", err)
	}
	ctx := context.Background()
	filter := additionalFilters
	if filter == nil {
		filter = map[string]interface{}{}
	}
	filter["_id"] = unhexID
	res := c.collection.FindOne(ctx, filter)
	if res == nil {
		return nil, custom_error.MakeErrorf("Not found. Result is nil.")
	}
	value, customErr := decoder(res.Decode)
	if customErr == nil {
		return value, nil
	}
	return nil, custom_error.NewErrorf(customErr, "Failed to decode. Collname: %v. Db-name: %v", c.collName, c.dbName)
}

func (c *collection) FindAll(decoder store.Decoder, filter map[string]interface{}, processor func(res interface{}) (bool, custom_error.CustomError)) custom_error.CustomError {
	ctx := context.Background()
	cursor, err := c.collection.Find(ctx, filter)
	if err != nil {
		return custom_error.MakeErrorf("Failed to find all. Error: %v", err)
	}
	if cursor == nil {
		return custom_error.MakeErrorf("Something went wrong. Empty cursor.")
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		value, customErr := decoder(cursor.Decode)
		if customErr != nil {
			return custom_error.NewErrorf(customErr, "Failed to decode.")
		}
		next, customErr := processor(value)
		if customErr != nil {
			return custom_error.NewErrorf(customErr, "Process failed.")
		}
		if !next {
			break
		}
	}
	err = cursor.Err()
	if err != nil {
		return custom_error.MakeErrorf("Failed to enum. Cursor error: %v", err)
	}
	return nil
}

func (c *collection) Create(object interface{}) (resObjID string, resErr custom_error.CustomError) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		resObjID = ""
		customErr, ok := r.(custom_error.CustomError)
		if ok {
			resErr = custom_error.NewErrorf(customErr, "Failed to create. Panic caught.")
			return
		}
		err, ok := r.(error)
		if ok {
			resErr = custom_error.MakeErrorf("Failed to create. Panic caught. Error: %v", err)
			return
		}
		resErr = custom_error.MakeErrorf("Unknown error. Panic caught. Error: %v", r)
	}()
	ctx := context.Background()
	res, err := c.collection.InsertOne(ctx, object)
	if err != nil {
		return "", custom_error.MakeErrorf("Failed to insert. Error: %v", err)
	}
	if res == nil {
		return "", custom_error.MakeErrorf("Failed to insert. Empty result.")
	}
	el, ok := res.InsertedID.(*bson.Element)
	if !ok {
		return "", custom_error.MakeErrorf("Inserted object has non-object-id type for ID. Unexpected type: %T", res.InsertedID)
	}
	val := el.Value()
	if val == nil {
		return "", custom_error.MakeErrorf("Something went wrong. Object ID is nil.")
	}
	return val.ObjectID().Hex(), nil
}

func (c *collection) Delete(id string, additionalFilters map[string]interface{}) (bool, custom_error.CustomError) {
	unhexID, err := objectid.FromHex(id)
	if err != nil {
		return false, custom_error.MakeErrorf("Failed to convert id. Error: %v", err)
	}
	ctx := context.Background()

	filter := additionalFilters
	if filter == nil {
		filter = map[string]interface{}{}
	}
	filter["_id"] = unhexID
	res, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, custom_error.MakeErrorf("Failed to remove. Error: %v", err)
	}
	if res == nil {
		return false, custom_error.MakeErrorf("Something went wrong. Empty result.")
	}

	return res.DeletedCount > 0, nil
}

func NewCollection(client *mgo.Client, dbName string, collectionName string) (store.Collection, custom_error.CustomError) {
	db := client.Database(dbName)
	if db == nil {
		return nil, custom_error.MakeErrorf("DB is nil")
	}
	coll := db.Collection(collectionName)
	if coll == nil {
		return nil, custom_error.MakeErrorf("Collection is nil")
	}
	return &collection{
		dbName:     dbName,
		collName:   collectionName,
		collection: coll,
	}, nil
}
