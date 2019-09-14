package resolver

import (
	"context"
	"log"

	"github.com/astrohot/backend/internal/api/model/user"
	"google.golang.org/api/iterator"
)

type queryResolver struct {
	*Resolver
}

func (r *queryResolver) Users(ctx context.Context) (users []*user.User, err error) {
	coll := r.Firestore.Client.Collection("users")
	iter := coll.Documents(ctx)
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err != nil {
			if err != iterator.Done {
				log.Println(err)
			}

			break
		}

		data := snapshot.Data()
		users = append(users, &user.User{
			ID:    snapshot.Ref.ID,
			Name:  data["name"].(string),
			Email: data["email"].(string),
			Likes: data["likes"].([]string),
		})
	}

	if err != nil && err != iterator.Done {
		return
	}

	return
}

func (r *queryResolver) Likes(ctx context.Context, mainID string) (likes []*string, err error) {
	coll := r.Firestore.Client.Collection("likes")
	query := coll.Where("mainID", "==", mainID)
	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err != nil {
			if err != iterator.Done {
				log.Println(err)
			}

			break
		}

		data := snapshot.Data()
		crushID := data["crushID"].(string)
		likes = append(likes, &crushID)
	}

	if err != nil && err != iterator.Done {
		likes = nil
	}

	return
}

func (r *queryResolver) Matches(ctx context.Context, mainID string) (matches []*string, err error) {
	coll := r.Firestore.Client.Collection("likes")
	query := coll.Where("mainID", "==", mainID)
	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err != nil {
			if err != iterator.Done {
				log.Println(err)
			}

			break
		}

		data := snapshot.Data()
		crushID := data["crushID"].(string)
		query = coll.Where("mainID", "==", crushID).Where("crushID", "==", mainID)
		matchesIter := query.Documents(ctx)

		// That query must return exactly zero or one values. If there's one
		// value, then it's a match ;D
		if _, err := matchesIter.Next(); err == nil {
			matches = append(matches, &crushID)
		}

		matchesIter.Stop()
	}

	if err != nil && err != iterator.Done {
		matches = nil
	}

	return
}
