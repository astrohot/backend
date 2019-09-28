package resolver

import (
	"context"

	"github.com/astrohot/backend/internal/api/generated"
	"github.com/astrohot/backend/internal/lib/auth"
	"github.com/astrohot/backend/internal/model/action"
	"github.com/astrohot/backend/internal/model/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type queryResolver struct {
	*Resolver
}

func (r *queryResolver) ValidateToken(ctx context.Context, token string) (bool, error) {
	_, err := auth.Parse(token)
	switch err {
	case nil:
		return true, nil
	case auth.ErrInvalidToken:
		return false, nil
	default:
		return false, err
	}
}

func (r *queryResolver) Auth(ctx context.Context, input generated.Auth) (*user.User, error) {
	var (
		u   user.User
		err error
	)

	// Check if user is already authenticated and thus it is inside the
	// context.
	u, ok := auth.FromContext(ctx).(user.User)

	// If not, then try to retrieve user from database
	if !ok || u.Token.Value == "" {
		u = u.AddFilter("email", input.Email)
		u, err = u.FindOne(ctx)

		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		// Try to authenticate user.
		var tok string
		tok, err = u.Authenticate(input.Password)
		if err != nil {
			return nil, err
		}

		u.Token = user.Token{
			Value:   tok,
			IsValid: true,
		}
	}

	return &u, nil
}

func (r *queryResolver) Users(ctx context.Context, offset, limit int) ([]*user.User, error) {
	// Check if user is authenticated.
	u, ok := auth.FromContext(ctx).(user.User)
	if !ok {
		return nil, ErrNotLogged
	}

	// Get list of users.
	u = u.AddFilter("_id", bson.M{"$ne": u.ID})
	us, err := u.Find(ctx)
	if err != nil {
		return nil, err
	}

	// Get all likes and dislikes (actions) of that user.
	as, err := action.Action{}.AddFilter("mainID", u.ID).Find(ctx)
	if err != nil {
		return nil, err
	}

	// Filter list of users excluding users with like or dislike.
	actionMap := make(map[primitive.ObjectID]struct{}, len(as))
	for _, a := range as {
		actionMap[a.CrushID] = struct{}{}
	}

	var list []*user.User
	for _, u := range us {
		if _, ok := actionMap[u.ID]; !ok {
			list = append(list, u)
		}
	}

	if offset < 0 || offset >= len(list) {
		offset = 0
	}

	if limit < 0 || limit > len(list) {
		limit = len(list)
	}

	return list[offset : limit+offset], nil
}

func (r *queryResolver) Matches(ctx context.Context, mainID primitive.ObjectID) ([]*primitive.ObjectID, error) {
	return nil, nil
}
