package resolver

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/astrohot/backend/internal/api/generated"
	"github.com/astrohot/backend/internal/database"
)

// Resolver ...
type Resolver struct {
	DB *database.DB
}

// Mutation ...
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

// Query ...
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}
