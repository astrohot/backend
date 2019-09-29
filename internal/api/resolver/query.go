package resolver

import (
	"context"
	"log"
	"sync"

	"github.com/astrohot/backend/internal/api/generated"
	"github.com/astrohot/backend/internal/domain/action"
	"github.com/astrohot/backend/internal/domain/user"
	"github.com/astrohot/backend/internal/lib/auth"
	"github.com/astrohot/backend/internal/lib/database"
	"github.com/astrohot/backend/internal/lib/zodiac"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type queryResolver struct {
	*Resolver
}

func (r *queryResolver) ValidateToken(ctx context.Context, token string) (bool, error) {
	_, err := auth.Parse(token)
	switch err {
	case nil:
		return true, nil
	case auth.ErrInvalidToken:
		return false, nil
	default:
		log.Println(err)
		return false, err
	}
}

func (r *queryResolver) Auth(ctx context.Context, input generated.Auth) (*user.User, error) {
	var (
		u   user.User
		err error
	)

	// Check if user is already authenticated and thus it is inside the
	// context.
	u, ok := auth.FromContext(ctx).(user.User)
	// If not, then try to retrieve user from database
	if !ok || u.Token.Value == "" {
		u = u.AddFilter("email", input.Email)
		u, err = u.FindOne(ctx)

		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		// Try to authenticate user.
		var tok string
		tok, err = u.Authenticate(input.Password)
		if err != nil {
			return nil, err
		}

		u.Token = user.Token{
			Value:   tok,
			IsValid: true,
		}
	}

	return &u, nil
}

func (r *queryResolver) GetUsers(ctx context.Context, offset, limit int) ([]*user.User, error) {
	// Check if user is authenticated.
	u, ok := auth.FromContext(ctx).(user.User)
	if !ok {
		return nil, ErrNotLogged
	}

	// Get list of users.
	u = u.AddFilter("_id", database.Filter{
		Field: "$ne",
		Value: u.ID,
	})

	us, err := u.Find(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Get all likes and dislikes (actions) of that user.
	as, err := action.Action{}.AddFilter("mainID", u.ID).Find(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Filter list of users excluding users with like or dislike.
	actionMap := make(map[primitive.ObjectID]struct{}, len(as))
	for _, a := range as {
		actionMap[a.CrushID] = struct{}{}
	}

	var list []*user.User
	for _, u := range us {
		if _, ok := actionMap[u.ID]; !ok {
			list = append(list, u)
		}
	}

	if offset < 0 || offset >= len(list) {
		offset = 0
	}

	if limit < 0 || limit > len(list) {
		limit = len(list)
	}

	return list[offset : limit+offset], nil
}

func (r *queryResolver) GetMatches(ctx context.Context, mainID primitive.ObjectID) ([]*primitive.ObjectID, error) {
	// Check if user is authenticated.
	u, ok := auth.FromContext(ctx).(user.User)
	if !ok {
		return nil, ErrNotLogged
	}

	var (
		wg        sync.WaitGroup
		likesFrom []*action.Action
		likesTo   []*action.Action
	)

	errCh := make(chan error, 2)
	quitCh := make(chan bool)

	// Handling errors.
	go func() {
		wg.Wait()

		close(errCh)
		shouldQuit := false

		for err := range errCh {
			if err != nil {
				log.Println(err)
				shouldQuit = true
			}
		}

		quitCh <- shouldQuit
	}()

	wg.Add(2)

	// Get Likes where mainID is u.ID
	go func() {
		var err error
		defer wg.Done()

		likesFrom, err = action.Action{}.AddFilter("$and", database.FilterList{
			{Field: "mainID", Value: u.ID},
			{Field: "type", Value: action.Like},
		}).Find(ctx)

		errCh <- err
	}()

	// Get Likes where u.ID is the crushID (reverse order).
	go func() {
		var err error
		defer wg.Done()

		likesTo, err = action.Action{}.AddFilter("$and", database.FilterList{
			{Field: "crushID", Value: u.ID},
			{Field: "type", Value: action.Like},
		}).Find(ctx)

		errCh <- err
	}()

	if shouldQuit := <-quitCh; shouldQuit {
		return nil, ErrFetchFailed
	}

	// Filter users to get only those where there're likes in both orders.
	isBoth := make(map[primitive.ObjectID]struct{}, len(likesTo))
	for _, l := range likesTo {
		isBoth[l.MainID] = struct{}{}
	}

	matches := []*primitive.ObjectID{}
	for _, l := range likesFrom {
		if _, ok := isBoth[l.CrushID]; ok {
			matches = append(matches, &l.CrushID)
		}
	}

	return matches, nil
}

func (r *queryResolver) GetHoroscope(ctx context.Context) (string, error) {
	// Check if user is authenticated.
	u, ok := auth.FromContext(ctx).(user.User)
	if !ok {
		return "", ErrNotLogged
	}

	u, err := u.AddFilter("_id", u.ID).FindOne(ctx)
	if err != nil {
		log.Println(err)
		return "", err
	}

	horoscope, err := zodiac.GetHoroscope(u.Sign)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return horoscope, err
}
