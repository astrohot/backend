package user

import (
	"time"

	"github.com/astrohot/backend/internal/lib/auth"
	"github.com/astrohot/backend/internal/lib/database"
	"github.com/astrohot/backend/internal/lib/zodiac"
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
	Name     string
	Email    string
	Password string
	Birth    time.Time
}

// User ...
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Token     Token              `bson:"-" json:"-"`
	Email     string             `bson:"email" json:"email,omitempty"`
	Password  string             `bson:"password" json:"password,omitempty"`
	Name      string             `bson:"name" json:"name,omitempty"`
	Sign      zodiac.Sign        `bson:"sign" json:"sign,omitempty"`
	Birth     time.Time          `bson:"birth" json:"birth,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt,omitempty"`
	Likes     []string           `bson:"-" json:"-"`
	where     database.FilterList
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

// ClearFilter ...
func (u User) ClearFilter() User {
	u.where = nil
	return u
}

// AddFilter ...
func (u User) AddFilter(field string, value interface{}) User {
	u.where = append(u.where, database.Filter{
		Field: field,
		Value: value,
	})

	return u
}

// Collection ...
func (u User) Collection() string {
	return "users"
}

// Where ...
func (u User) Where() database.FilterList {
	return u.where
}
