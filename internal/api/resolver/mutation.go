package resolver

import (
	"context"
	"log"

	"github.com/astrohot/backend/internal/domain/action"
	"github.com/astrohot/backend/internal/domain/user"
	"github.com/astrohot/backend/internal/lib/auth"
	"github.com/astrohot/backend/internal/lib/zodiac"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) CreateUser(ctx context.Context, input user.NewUser) (*user.User, error) {
	// Check if user already exists. If it exists then err must be
	// mongo.ErrNoDocuments. Otherwise, err is nil.
	u := user.User{}.AddFilter("email", input.Email)
	if _, err := u.FindOne(ctx); err == nil {
		return nil, ErrUserExists
	}

	password, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.MinCost,
	)

	if err != nil {
		return nil, err
	}

	u = user.User{
		Email:    input.Email,
		Password: string(password),
		Name:     input.Name,
		Birth:    input.Birth,
		Sign:     zodiac.GetSign(input.Birth),
	}

	if u, err = u.Insert(ctx); err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

func (r *mutationResolver) CreateLike(ctx context.Context, crushID primitive.ObjectID) (*action.Action, error) {
	// Check if user is authenticated. If it's not, return with error.
	u, _ := auth.FromContext(ctx).(user.User)
	if u.Token.Value == "" {
		return nil, ErrNotLogged

	}

	// Check if u.ID is not the same as crushID.
	if u.ID == crushID {
		return nil, ErrMainEqualsCrush

	}

	a := action.Action{
		MainID:  u.ID,
		CrushID: crushID,
		Type:    action.Like,
	}

	// Create the new like.
	a, err := a.Insert(ctx)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *mutationResolver) CreateDislike(ctx context.Context, crushID primitive.ObjectID) (*action.Action, error) {
	// Check if user is authenticated. If it's not, return with error.
	u, _ := auth.FromContext(ctx).(user.User)
	if u.Token.Value == "" {
		return nil, ErrNotLogged

	}

	// Check if u.ID is not the same as crushID.
	if u.ID == crushID {
		return nil, ErrMainEqualsCrush

	}

	a := action.Action{
		MainID:  u.ID,
		CrushID: crushID,
		Type:    action.Dislike,
	}

	// Create the new like.
	a, err := a.Insert(ctx)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
