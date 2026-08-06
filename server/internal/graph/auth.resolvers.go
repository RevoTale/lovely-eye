package graph

import (
	"context"
	"fmt"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/graph/model"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthPayload, error) {
	user, tokens, err := r.AuthService.Register(ctx, auth.RegisterInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	if w := GetResponseWriter(ctx); w != nil {
		r.AuthCookies.SetAuthCookies(w, tokens)
	}

	return &model.AuthPayload{
		User: convertToGraphQLUser(user),
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	user, tokens, err := r.AuthService.Login(ctx, auth.LoginInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	if w := GetResponseWriter(ctx); w != nil {
		r.AuthCookies.SetAuthCookies(w, tokens)
	}

	return &model.AuthPayload{
		User: convertToGraphQLUser(user),
	}, nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	if w := GetResponseWriter(ctx); w != nil {
		r.AuthCookies.ClearAuthCookies(w)
	}
	return true, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, nil //nolint:nilnil // unauthenticated returns no user without error
	}

	user, err := r.AuthService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return convertToGraphQLUser(user), nil
}

// RegistrationStatus is the resolver for the registrationStatus field.
func (r *queryResolver) RegistrationStatus(ctx context.Context) (*model.RegistrationStatus, error) {
	status, err := r.AuthService.RegistrationStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get registration status: %w", err)
	}

	return &model.RegistrationStatus{
		HasUsers:          status.HasUsers,
		AllowRegistration: status.AllowRegistration,
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
