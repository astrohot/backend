package user

import (
	"time"

	"github.com/astrohot/backend/internal/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Token ...
type Token struct {
	Value   string
	IsValid bool
}

// NewUser ...
type NewUser struct {
	Name        string
	Email       string
	Password    string
	Description string
	Birth       string
}

// User ...
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Token       Token              `bson:"-" json:"-"`
	Email       string             `bson:"email" json:"email,omitempty"`
	Password    string             `bson:"password" json:"password,omitempty"`
	Name        string             `bson:"name" json:"name,omitempty"`
	Description string             `bson:"description" json:"description,omitempty"`
	Birth       string             `bson:"birth" json:"birth,omitempty"`
	Likes       []string           `bson:"-" json:"-"`
}

// Authenticate returns a valid token string or an error. It implements the
// auth.Auther interface.
func (u User) Authenticate(plain string) (string, error) {
	if u.Token.Value != "" && u.Token.IsValid {
		return u.Token.Value, nil
	}

	bytePlain := []byte(plain)
	byteHash := []byte(u.Password)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return "", err
	}

	// Generate a new token.
	tok, err := auth.CreateToken(u.Email, time.Now().Add(time.Hour*24))
	if err != nil {
		return "", err
	}

	return tok, nil
}
