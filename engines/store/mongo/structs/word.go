package structs

import "github.com/coldze/mongo-go-driver/bson/objectid"

type Word struct {
	ID           objectid.ObjectID `bson:"_id"`
	ClientID     objectid.ObjectID `bson:"client_id"`
	BoardID      objectid.ObjectID `bson:"board_id"`
	Text         string            `bson:"text"`
	Translations []string          `bson:"translations"`
}
