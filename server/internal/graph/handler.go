package graph

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
)

type contextKey string

const (
	responseWriterKey contextKey = "response_writer"
)

// Handler creates a GraphQL handler using gqlgen's generated ExecutableSchema
func Handler(resolver *Resolver) http.HandlerFunc {
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers: resolver,
	}))

	return func(w http.ResponseWriter, r *http.Request) {
		// Add ResponseWriter to context for cookie setting in resolvers
		ctx := context.WithValue(r.Context(), responseWriterKey, w)
		srv.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetResponseWriter retrieves the ResponseWriter from context.
func GetResponseWriter(ctx context.Context) http.ResponseWriter {
	w, ok := ctx.Value(responseWriterKey).(http.ResponseWriter)
	if !ok {
		return nil
	}
	return w
}
