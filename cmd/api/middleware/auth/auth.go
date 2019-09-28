package auth

import (
	"net/http"

	"github.com/astrohot/backend/internal/lib/auth"
	"github.com/astrohot/backend/internal/model/user"
)

const authHeader = "Authorization"

// Middleware decodes the authorization header and insert user into the context.
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get(authHeader)
			u := user.User{}

			// If token is empty go ahead without inserting the user into the
			// context.
			if tokenString == "" {
				ctx := auth.WithContext(r.Context(), u)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Parse token string and get user email from it.
			token, err := auth.Parse(tokenString)
			switch err {
			case nil:
				u = u.AddFilter("email", token.Email)
				u, err := u.FindOne(r.Context())

				if err != nil {
					http.Error(w, "user not found", http.StatusUnauthorized)
					return
				}

				// Clear filter for future uses.
				u = u.ClearFilter()

				// Put user into the context.
				u.Token = user.Token{
					Value:   tokenString,
					IsValid: true,
				}

				ctx := auth.WithContext(r.Context(), u)

				// Call the next with our new context.
				next.ServeHTTP(w, r.WithContext(ctx))

			case auth.ErrInvalidToken:
				u := user.User{
					Token: user.Token{
						Value:   tokenString,
						IsValid: false,
					},
				}

				ctx := auth.WithContext(r.Context(), u)

				// Call the next with our new context.
				next.ServeHTTP(w, r.WithContext(ctx))

			default:
				http.Error(w, "Error while parsing jwt token string", http.StatusUnauthorized)
				return
			}
		})
	}
}
