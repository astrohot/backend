package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/astrohot/backend/internal/api/generated"
	"github.com/astrohot/backend/internal/api/resolver"
	"github.com/astrohot/backend/internal/database"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	firestore, err := database.New(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	defer firestore.Client.Close()

	http.Handle("/", handler.Playground("GraphQL playground", "/graphql"))
	http.Handle("/graphql", handler.GraphQL(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &resolver.Resolver{
					Firestore: firestore,
				},
			},
		)))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}