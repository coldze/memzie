package mongo

import (
	"context"

	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/mongo-go-driver/bson"
	"github.com/coldze/mongo-go-driver/core/options"
	mgo "github.com/coldze/mongo-go-driver/mongo"
	"github.com/coldze/primitives/custom_error"
)

type collectionWrap struct {
	dbName     string
	collName   string
	collection *mgo.Collection
}

func (c *collectionWrap) findOne(decoder store.Decoder, filter map[string]interface{}, opts ...options.FindOneOptioner) (interface{}, custom_error.CustomError) {
	ctx := context.Background()
	res := c.collection.FindOne(ctx, filter, opts...)
	if res == nil {
		return nil, custom_error.MakeErrorf("Not found. Result is nil.")
	}
	value, customErr := decoder(res.Decode)
	if customErr == nil {
		return value, nil
	}
	return nil, custom_error.NewErrorf(customErr, "Failed to decode. Collname: %v. Db-name: %v", c.collName, c.dbName)
}

func (c *collectionWrap) findAll(decoder store.Decoder, processor func(res interface{}) (bool, custom_error.CustomError), filter map[string]interface{}, opts ...options.FindOptioner) custom_error.CustomError {
	ctx := context.Background()
	cursor, err := c.collection.Find(ctx, filter, opts...)
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

func (c *collectionWrap) create(object interface{}, opts ...options.InsertOneOptioner) (resObjID string, resErr custom_error.CustomError) {
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
	res, err := c.collection.InsertOne(ctx, object, opts...)
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

func (c *collectionWrap) delete(filter map[string]interface{}, opts ...options.DeleteOptioner) (bool, custom_error.CustomError) {
	ctx := context.Background()
	res, err := c.collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return false, custom_error.MakeErrorf("Failed to remove. Error: %v", err)
	}
	if res == nil {
		return false, custom_error.MakeErrorf("Something went wrong. Empty result.")
	}

	return res.DeletedCount > 0, nil
}

type collection struct {
	underlying *collectionWrap
}

func (c *collection) FindOne(decoder store.Decoder, filter map[string]interface{}) (interface{}, custom_error.CustomError) {
	res, customErr := c.underlying.findOne(decoder, filter)
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to findOne.")
	}
	return res, nil
}

func (c *collection) FindAll(decoder store.Decoder, filter map[string]interface{}, processor func(res interface{}) (bool, custom_error.CustomError)) custom_error.CustomError {
	customErr := c.underlying.findAll(decoder, processor, filter)
	if customErr != nil {
		return custom_error.NewErrorf(customErr, "Failed to find all.")
	}
	return nil
}

func (c *collection) Create(object interface{}) (resObjID string, resErr custom_error.CustomError) {
	res, err := c.underlying.create(object)
	if err != nil {
		return "", custom_error.NewErrorf(err, "Failed to insert.")
	}
	return res, nil
}

func (c *collection) Delete(filter map[string]interface{}) (bool, custom_error.CustomError) {
	res, customErr := c.underlying.delete(filter)
	if customErr != nil {
		return false, custom_error.NewErrorf(customErr, "Failed to remove.")
	}
	return res, nil
}

type logicCollection struct {
	underlying *collectionWrap
	sort       *bson.Document
}

func (c *logicCollection) FindOne(decoder store.Decoder, filter map[string]interface{}) (interface{}, custom_error.CustomError) {
	res, customErr := c.underlying.findOne(decoder, filter, options.OptSort{
		Sort: c.sort,
	})
	if customErr != nil {
		return nil, custom_error.NewErrorf(customErr, "Failed to findOne.")
	}
	return res, nil
}

func (c *logicCollection) FindAll(decoder store.Decoder, filter map[string]interface{}, processor func(res interface{}) (bool, custom_error.CustomError)) custom_error.CustomError {
	return custom_error.MakeErrorf("Should not be called.")
}

func (c *logicCollection) Create(object interface{}) (resObjID string, resErr custom_error.CustomError) {
	return "", custom_error.MakeErrorf("Should not be called.")
}

func (c *logicCollection) Delete(id string, additionalFilters map[string]interface{}) (bool, custom_error.CustomError) {
	return false, custom_error.MakeErrorf("Should not be called.")
}

func NewCollectionFactory(client *mgo.Client, dbName string) func(collName string) (store.Collection, custom_error.CustomError) {
	return func(collectionName string) (store.Collection, custom_error.CustomError) {
		db := client.Database(dbName)
		if db == nil {
			return nil, custom_error.MakeErrorf("DB is nil")
		}
		coll := db.Collection(collectionName)
		if coll == nil {
			return nil, custom_error.MakeErrorf("Collection is nil")
		}
		return &collection{
			underlying: &collectionWrap{
				dbName:     dbName,
				collName:   collectionName,
				collection: coll,
			},
		}, nil
	}
}

func NewLogicCollection(client *mgo.Client, dbName string, collectionName string, sortOrder map[string]interface{}) (store.Collection, custom_error.CustomError) {
	/*sortDoc*/ _, err := mgo.TransformDocument(sortOrder)
	if err != nil {
		return nil, custom_error.MakeErrorf("Failed to convert sort order. Error: %v", err)
	}
	db := client.Database(dbName)
	if db == nil {
		return nil, custom_error.MakeErrorf("DB is nil")
	}
	coll := db.Collection(collectionName)
	if coll == nil {
		return nil, custom_error.MakeErrorf("Collection is nil")
	}
	/*return &logicCollection{
		underlying: &collectionWrap{
			dbName:     dbName,
			collName:   collectionName,
			collection: coll,
			sort:       sortDoc,
		},
	}, nil*/
	return nil, custom_error.MakeErrorf("Not implemented.")
}
