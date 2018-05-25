package structs

import "github.com/mongodb/mongo-go-driver/bson/objectid"

type Folder struct {
	ID          objectid.ObjectID `bson:"_id"`
	ClientID    objectid.ObjectID `bson:"client_id"`
	Name        string            `bson:"Name"`
	Description string            `bson:"Description"`
}
