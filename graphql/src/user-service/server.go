package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"user-service/graph"
	"user-service/graph/model"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := connectToPostgres()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Db: db,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connectToPostgres() *gorm.DB {
	server := os.Getenv("POSTGRES_HOST")
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		panic("Could not parse PORT to an integar")
	}
	name := os.Getenv("POSTGRES_DATABASE")
	user := os.Getenv("POSTGRES_USERNAME")
	pass := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%d", server, user, pass, name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to postgres: %v", err)
	}

	db.AutoMigrate(
		&model.User{},
	)

	return db
}
