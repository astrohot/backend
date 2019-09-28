package action

import (
	"context"

	"github.com/astrohot/backend/internal/lib/database"
)

// Insert ...
func (a Action) Insert(ctx context.Context) (Action, error) {
	insertedID, err := database.InsertOne(ctx, a)
	if err != nil {
		return Action{}, err
	}

	a.ID = insertedID
	return a, nil
}

// Find ...
func (a Action) Find(ctx context.Context) (as []*Action, err error) {
	cursor, err := database.Find(ctx, a)
	if err != nil {
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result Action
		if err = cursor.Decode(&result); err != nil {
			as = nil
			break
		}

		as = append(as, &result)
	}

	return
}
