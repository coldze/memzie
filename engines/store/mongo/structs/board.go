package structs

import "github.com/mongodb/mongo-go-driver/bson/objectid"

type Board struct {
	ID          objectid.ObjectID `bson:"_id"`
	ClientID    objectid.ObjectID `bson:"client_id"`
	Name        string            `bson:"Name"`
	Description string            `bson:"Description"`
}
