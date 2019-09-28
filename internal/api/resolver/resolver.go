package resolver

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/astrohot/backend/internal/api/generated"
)

// Resolver ...
type Resolver struct{}

// Mutation ...
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

// Query ...
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// User ...
func (r *Resolver) User() generated.UserResolver {
	return &userResolver{r}
}
