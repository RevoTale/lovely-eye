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

func Handler(resolver *Resolver) http.HandlerFunc {
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers: resolver,
	}))

	return func(w http.ResponseWriter, r *http.Request) {

		ctx := context.WithValue(r.Context(), responseWriterKey, w)
		srv.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetResponseWriter(ctx context.Context) http.ResponseWriter {
	w, ok := ctx.Value(responseWriterKey).(http.ResponseWriter)
	if !ok {
		return nil
	}
	return w
}
