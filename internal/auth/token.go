package auth

import (
	"errors"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// JWT key used to create the signature.
	jwtKey = []byte("my_secret_key")

	// ErrInvalidToken ...
	ErrInvalidToken = errors.New("invalid token")
)

// Claims it's a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Email string `bson:"email" json:"email"`
	jwt.StandardClaims
}

// CreateToken ...
func CreateToken(email string, expiration time.Time) (string, error) {
	claims := Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expiration.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Parse the JWT string and store the result in `claims`.
// Note that we are passing the key in this method as well. This method will return an error
// if the token is invalid (if it has expired according to the expiry time we set on sign in),
// or if the signature does not match
func Parse(tokenString string) (*Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}
