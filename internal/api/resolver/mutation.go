package resolver

import (
	"context"
	"log"
	"time"

	"github.com/astrohot/backend/internal/auth"
	"github.com/astrohot/backend/internal/model/like"
	"github.com/astrohot/backend/internal/model/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) CreateUser(ctx context.Context, input user.NewUser) (*user.User, error) {
	// First insert the credentials into user collection
	coll := r.DB.Collection("users")

	// Check if user already exists
	exists := r.DB.Exists(ctx, coll, bson.M{"email": input.Email})
	if exists {
		return nil, ErrUserExists
	}

	password, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.MinCost,
	)

	if err != nil {
		return nil, err
	}

	u := user.User{
		Email:       input.Email,
		Password:    string(password),
		Name:        input.Name,
		Description: input.Description,
		Birth:       input.Birth,
	}

	result, err := coll.InsertOne(ctx, u)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	u.ID = result.InsertedID.(primitive.ObjectID)
	tok, err := auth.CreateToken(u.Email, time.Now().Add(time.Hour*24))
	if err != nil {
		return nil, err
	}

	u.Token = user.Token{
		Value:   tok,
		IsValid: true,
	}

	return &u, nil
}

func (r *mutationResolver) CreateLike(ctx context.Context, input like.NewLike) (*like.Like, error) {
	// Check if user is authenticated. If it's not, return with error.
	u, _ := auth.FromContext(ctx).(user.User)
	if u.Token.Value == "" {
		return nil, ErrNotLogged

	}

	// Check if input.MainID refers to the user ID.
	if u.ID != input.MainID {
		return nil, ErrMismatchMainID
	}

	// Check if input.MainID is not the same as input.CrushID.
	if input.MainID == input.CrushID {
		return nil, ErrMainEqualsCrush

	}

	// Create the new like.
	coll := r.DB.Collection("likes")
	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	l := like.Like{
		ID:      result.InsertedID.(primitive.ObjectID),
		MainID:  u.ID,
		CrushID: input.CrushID,
	}

	return &l, nil
}
