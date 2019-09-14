package resolver

import (
	"context"

	"github.com/astrohot/backend/internal/api/model/like"
	"github.com/astrohot/backend/internal/api/model/user"
)

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) CreateUser(ctx context.Context, input user.NewUser) (entity *user.User, err error) {
	coll := r.Firestore.Client.Collection("users")
	document, _, err := coll.Add(ctx, input)
	if err != nil {
		return
	}

	entity = &user.User{
		ID:    document.ID,
		Name:  input.Name,
		Email: input.Email,
	}

	return
}

func (r *mutationResolver) CreateLike(ctx context.Context, input like.NewLike) (entity *like.Like, err error) {
	coll := r.Firestore.Client.Collection("likes")
	document, _, err := coll.Add(ctx, input)
	if err != nil {
		return
	}

	entity = &like.Like{
		ID:      document.ID,
		MainID:  input.MainID,
		CrushID: input.CrushID,
	}

	return
}
