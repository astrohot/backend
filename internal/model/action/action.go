package action

import (
	"github.com/astrohot/backend/internal/lib/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Action types
const (
	Like    = "like"
	Dislike = "dislike"
)

// NewAction ...
type NewAction struct {
	MainID  primitive.ObjectID
	CrushID primitive.ObjectID
}

// Action ...
type Action struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	MainID  primitive.ObjectID `bson:"mainID" json:"mainID,omitempty"`
	CrushID primitive.ObjectID `bson:"crushID" json:"crushID,omitempty"`
	Type    string             `bson:"action" json:"action,omitempty"`
	where   database.FilterList
}

// AddFilter ...
func (a Action) AddFilter(field string, value interface{}) Action {
	a.where = append(a.where, database.Filter{
		Field: field,
		Value: value,
	})

	return a
}

// Collection ...
func (a Action) Collection() string {
	return "actions"
}

// Where ...
func (a Action) Where() database.FilterList {
	return a.where
}
