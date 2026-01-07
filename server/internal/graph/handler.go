package graph

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
)

// Handler creates a GraphQL handler using gqlgen's generated ExecutableSchema
func Handler(resolver *Resolver) http.HandlerFunc {
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers: resolver,
	}))

	return srv.ServeHTTP
}
