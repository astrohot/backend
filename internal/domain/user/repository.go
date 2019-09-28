package user

import (
	"context"
	"time"

	"github.com/astrohot/backend/internal/lib/auth"
	"github.com/astrohot/backend/internal/lib/database"
)

// FindOne ...
func (u User) FindOne(ctx context.Context) (User, error) {
	err := database.FindOne(ctx, &u)
	return u, err
}

// Find ...
func (u User) Find(ctx context.Context) (us []*User, err error) {
	cursor, err := database.Find(ctx, u)
	if err != nil {
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result User
		if err = cursor.Decode(&result); err != nil {
			us = nil
			break
		}

		us = append(us, &result)
	}

	return
}

// Insert ...
func (u User) Insert(ctx context.Context) (User, error) {
	u.CreatedAt = time.Now().UTC()
	insertedID, err := database.InsertOne(ctx, u)
	if err != nil {
		return User{}, err
	}

	tok, err := auth.CreateToken(u.Email, time.Now().Add(time.Hour*24))
	if err != nil {
		return User{}, err
	}

	u.ID = insertedID
	u.Token = Token{
		Value:   tok,
		IsValid: true,
	}

	return u, nil
}
