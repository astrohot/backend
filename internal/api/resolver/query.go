package resolver

import (
	"context"
	"log"

	"github.com/astrohot/backend/internal/api/generated"
	"github.com/astrohot/backend/internal/auth"
	"github.com/astrohot/backend/internal/model/like"
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
		u, err = r.DB.GetUserByEmail(ctx, input.Email)
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

func (r *queryResolver) Users(ctx context.Context) (us []*user.User, err error) {
	return
}

func (r *queryResolver) Likes(ctx context.Context, mainID primitive.ObjectID) ([]*primitive.ObjectID, error) {
	// Check if user is authenticated. If it's not, return with error.
	u, ok := auth.FromContext(ctx).(user.User)
	if !ok {
		return nil, ErrNotLogged
	}

	// Get likes associated with that user.
	coll := r.DB.Collection("likes")
	where := bson.M{"mainID": u.ID}

	cursor, err := coll.Find(ctx, where)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ls []*primitive.ObjectID
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var l like.Like
		if err = cursor.Decode(&l); err != nil {
			log.Println(err)
			ls = nil
			break
		}

		ls = append(ls, &l.CrushID)
	}

	return ls, err
}

func (r *queryResolver) Matches(ctx context.Context, mainID primitive.ObjectID) ([]*primitive.ObjectID, error) {
	return nil, nil
}
