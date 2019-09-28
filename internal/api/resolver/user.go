package resolver

import (
	"context"

	"github.com/astrohot/backend/internal/domain/user"
)

type userResolver struct {
	*Resolver
}

func (r *userResolver) Sign(ctx context.Context, u *user.User) (string, error) {
	return u.Sign.String(), nil
}
