package resolver

import (
	"context"

	"github.com/astrohot/backend/internal/api/generated"
	"github.com/astrohot/backend/internal/api/model/user"
	"google.golang.org/api/iterator"
)

type queryResolver struct {
	*Resolver
}

func (r *queryResolver) Users(ctx context.Context) ([]*user.User, error) {
	users := []*user.User{}
	iter := r.Firestore.Client.Collection("Users").Documents(ctx)
	document, err := iter.Next()

	if err != nil && err != iterator.Done {
		return nil, err
	}

	for ; err != iterator.Done; document, err = iter.Next() {
		data := document.Data()
		users = append(users, &user.User{
			ID:    document.Ref.ID,
			Name:  data["Name"].(string),
			Email: data["Email"].(string),
		})
	}

	if err != nil && err != iterator.Done {
		return nil, err
	}

	return users, nil
}

func (r *queryResolver) Matches(ctx context.Context) ([]*generated.Match, error) {
	matches := []*generated.Match{}
	iter := r.Firestore.Client.Collection("Matches").Documents(ctx)
	document, err := iter.Next()

	if err != nil && err != iterator.Done {
		return nil, err
	}

	for ; err != iterator.Done; document, err = iter.Next() {
		data := document.Data()
		matches = append(matches, &generated.Match{
			UserA: data["UserA"].(string),
			UserB: data["UserB"].(string),
		})
	}

	if err != nil && err != iterator.Done {
		return nil, err
	}

	return matches, nil
}
