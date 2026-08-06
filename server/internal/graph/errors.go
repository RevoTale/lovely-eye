package graph

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/country"
	"github.com/lovely-eye/server/internal/event"
	"github.com/lovely-eye/server/internal/site"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type errorCode string

const (
	errorCodeBadUserInput    errorCode = "BAD_USER_INPUT"
	errorCodeConflict        errorCode = "CONFLICT"
	errorCodeForbidden       errorCode = "FORBIDDEN"
	errorCodeInternal        errorCode = "INTERNAL_SERVER_ERROR"
	errorCodeNotFound        errorCode = "NOT_FOUND"
	errorCodeUnauthenticated errorCode = "UNAUTHENTICATED"
)

type operationError struct {
	code    errorCode
	message string
}

func (e *operationError) Error() string { return e.message }

func badUserInput(message string) error {
	return &operationError{code: errorCodeBadUserInput, message: message}
}

func badUserInputf(format string, args ...any) error {
	return badUserInput(fmt.Sprintf(format, args...))
}

func unauthenticated() error {
	return &operationError{code: errorCodeUnauthenticated, message: "unauthorized"}
}

func presentError(ctx context.Context, err error) *gqlerror.Error {
	presented := graphql.DefaultErrorPresenter(ctx, err)
	code := classifyError(err)

	if presented.Extensions == nil {
		presented.Extensions = make(map[string]any, 1)
	}
	presented.Extensions["code"] = string(code)

	if code == errorCodeInternal {
		slog.ErrorContext(ctx, "GraphQL operation failed", "path", presented.Path.String(), "error", err)
		presented.Message = "internal server error"
	}

	return presented
}

func classifyError(err error) errorCode {
	var operationErr *operationError
	if errors.As(err, &operationErr) {
		return operationErr.code
	}

	switch {
	case errors.Is(err, site.ErrNotAuthorized), errors.Is(err, auth.ErrRegistrationDisabled):
		return errorCodeForbidden
	case errors.Is(err, site.ErrSiteNotFound), errors.Is(err, country.ErrNotFound):
		return errorCodeNotFound
	case errors.Is(err, site.ErrSiteExists), errors.Is(err, auth.ErrUserExists):
		return errorCodeConflict
	case errors.Is(err, auth.ErrUserNotFound):
		return errorCodeUnauthenticated
	case isValidationError(err):
		return errorCodeBadUserInput
	}

	var graphQLError *gqlerror.Error
	if errors.As(err, &graphQLError) && graphQLError.Err == nil {
		return errorCodeBadUserInput
	}

	return errorCodeInternal
}

func isValidationError(err error) bool {
	return errors.Is(err, auth.ErrInvalidCredentials) ||
		errors.Is(err, auth.ErrInvalidUsername) ||
		errors.Is(err, auth.ErrInvalidPassword) ||
		errors.Is(err, site.ErrInvalidDomain) ||
		errors.Is(err, site.ErrInvalidSiteName) ||
		errors.Is(err, site.ErrDomainTooLong) ||
		errors.Is(err, site.ErrSiteNameTooLong) ||
		errors.Is(err, site.ErrInvalidIPAddress) ||
		errors.Is(err, site.ErrInvalidCountryCode) ||
		errors.Is(err, site.ErrTooManyBlockedIPs) ||
		errors.Is(err, site.ErrTooManyBlockedCountries) ||
		errors.Is(err, event.ErrInvalidEventName) ||
		errors.Is(err, event.ErrInvalidFieldKey) ||
		errors.Is(err, event.ErrInvalidFieldType) ||
		errors.Is(err, event.ErrInvalidFieldLimit)
}
