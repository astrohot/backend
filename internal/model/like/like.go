package like

import "go.mongodb.org/mongo-driver/bson/primitive"

// Like ...
type Like struct {
	ID      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	MainID  primitive.ObjectID `bson:"mainID" json:"mainID,omitempty"`
	CrushID primitive.ObjectID `bson:"crushID" json:"crushID,omitempty"`
}

// NewLike ...
type NewLike struct {
	MainID  primitive.ObjectID `bson:"mainID" json:"mainID"`
	CrushID primitive.ObjectID `bson:"crushID" json:"crushID"`
}
