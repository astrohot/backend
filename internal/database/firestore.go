package database

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const defaultFCredentials = "firebase.json"

type Firestore struct {
	Client *firestore.Client
}

func New(ctx context.Context) (*Firestore, error) {
	fcredentials := os.Getenv("FCREDENTIALS")
	if fcredentials == "" {
		fcredentials = defaultFCredentials
	}

	opt := option.WithCredentialsFile(fcredentials)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	firestore := Firestore{
		Client: client,
	}

	return &firestore, nil
}
