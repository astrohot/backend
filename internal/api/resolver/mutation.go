package resolver

import (
	"context"

	"github.com/astrohot/backend/internal/api/generated"
	"github.com/astrohot/backend/internal/api/model/user"
)

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) CreateUser(ctx context.Context, input generated.NewUser) (*user.User, error) {
	document, _, err := r.Firestore.Client.Collection("Users").Add(ctx, input)
	if err != nil {
		return nil, err
	}

	user := user.User{
		ID:    document.ID,
		Name:  input.Name,
		Email: input.Email,
	}

	return &user, nil
}

func (r *mutationResolver) MatchUsers(ctx context.Context, userA, userB string) (*generated.Match, error) {
	match := generated.Match{
		UserA: userA,
		UserB: userB,
	}

	_, _, err := r.Firestore.Client.Collection("Matches").Add(ctx, match)
	if err != nil {
		return nil, err
	}

	return &match, nil
}
