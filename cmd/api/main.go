package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/astrohot/backend/cmd/api/middleware/auth"
	"github.com/astrohot/backend/internal/api/generated"
	"github.com/astrohot/backend/internal/api/resolver"
	"github.com/astrohot/backend/internal/lib/database"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("undefined database name")
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		log.Fatal("undefined database uri")
	}

	err := database.Init(context.Background(), dbName, dbURI)
	if err != nil {
		log.Fatalln(err)
	}

	router := chi.NewRouter()

	// Add middlewares around every request.
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Content-type", "Authorization"},
	}).Handler)
	router.Use(auth.Middleware())

	router.Handle("/", handler.Playground("GraphQL playground", "/graphql"))
	router.Handle("/graphql", handler.GraphQL(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: &resolver.Resolver{}},
		)),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
