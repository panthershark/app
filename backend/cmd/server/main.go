package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/panthershark/app/backend/internal/graph"
	"github.com/panthershark/app/backend/internal/graph/generated"
	"github.com/panthershark/app/backend/internal/reqctx"
	"github.com/rs/cors"
)

const (
	queryPath      string = "/query"
	playgroundPath string = "/playground"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	serverPort, err := strconv.ParseInt(os.Getenv("SERVER_PORT"), 10, 64)
	if err != nil {
		log.Panicf("bad port: %s", os.Getenv("SERVER_PORT"))
	}

	r := chi.NewRouter()

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		Debug:            true,
	}).Handler)

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(reqctx.GraphQLMiddleware)

	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	gqlPlayground := playground.Handler("GraphQL", queryPath)

	r.Mount(queryPath, gqlServer)
	r.Mount(playgroundPath, gqlPlayground)

	// Log the settings
	b := strings.Builder{}
	b.WriteString("Starting Server\n")
	b.WriteString(fmt.Sprintf("Allowed Origins: %v\n", strings.Join(allowedOrigins, " ")))
	b.WriteString(fmt.Sprintf("Server Port: %v\n", serverPort))
	b.WriteString(fmt.Sprintf("GraphQL: http://localhost:%d%s\n", serverPort, queryPath))
	b.WriteString(fmt.Sprintf("GraphQL Playground: http://localhost:%d%s\n", serverPort, playgroundPath))
	log.Println(b.String())

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), r))

}
