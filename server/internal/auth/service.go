package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserExists           = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrRegistrationDisabled = errors.New("registration is disabled")
	ErrInvalidUsername      = errors.New("invalid username")
	ErrInvalidPassword      = errors.New("invalid password")
)

const (
	accessTokenCookie  = "le_access"
	refreshTokenCookie = "le_refresh"
	minPasswordBytes   = 8
	maxPasswordBytes   = 72
	maxUsernameBytes   = 128
)

type Config struct {
	JWTSecret         string
	AccessTokenExpiry time.Duration
	RefreshExpiry     time.Duration
	AllowRegistration bool
	SecureCookies     bool
	CookieDomain      string
}

// jwtService implements the Service interface.
type jwtService struct {
	userRepo          *repository.UserRepository
	jwt               *jwtProvider
	allowRegistration bool
	secureCookies     bool
	cookieDomain      string
	accessExpiry      time.Duration
	refreshExpiry     time.Duration
}

// NewService creates a new authentication service.
func NewService(userRepo *repository.UserRepository, cfg Config) Service {
	return &jwtService{
		userRepo:          userRepo,
		jwt:               newJWTProvider(cfg.JWTSecret, cfg.AccessTokenExpiry, cfg.RefreshExpiry),
		allowRegistration: cfg.AllowRegistration,
		secureCookies:     cfg.SecureCookies,
		cookieDomain:      cfg.CookieDomain,
		accessExpiry:      cfg.AccessTokenExpiry,
		refreshExpiry:     cfg.RefreshExpiry,
	}
}

func (s *jwtService) Register(ctx context.Context, input RegisterInput) (*User, *Tokens, error) {
	username, err := normalizeUsername(input.Username)
	if err != nil {
		return nil, nil, err
	}
	if err := validatePassword(input.Password); err != nil {
		return nil, nil, err
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return nil, nil, err
	}

	dbUser := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.CreateForRegistration(ctx, dbUser, s.allowRegistration); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, nil, ErrUserExists
		}
		if errors.Is(err, repository.ErrUserRegistrationDisabled) {
			return nil, nil, ErrRegistrationDisabled
		}
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	user := &User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
		Role:     dbUser.Role,
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *jwtService) Login(ctx context.Context, input LoginInput) (*User, *Tokens, error) {
	username, err := normalizeUsername(input.Username)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	if err := validatePassword(input.Password); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	dbUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if !checkPassword(input.Password, dbUser.PasswordHash) {
		return nil, nil, ErrInvalidCredentials
	}

	user := &User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
		Role:     dbUser.Role,
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *jwtService) RefreshTokens(ctx context.Context, refreshToken string) (*Tokens, error) {
	claims, err := s.jwt.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	dbUser, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user := &User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
		Role:     dbUser.Role,
	}

	return s.generateTokens(user)
}

func (s *jwtService) ValidateAccessToken(token string) (*Claims, error) {
	return s.jwt.validateAccessToken(token)
}

func (s *jwtService) GetUserByID(ctx context.Context, id int64) (*User, error) {
	dbUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
		Role:     dbUser.Role,
	}, nil
}

func (s *jwtService) RegistrationStatus(ctx context.Context) (*RegistrationStatus, error) {
	isFirstUser, err := s.isFirstUser(ctx)
	if err != nil {
		return nil, err
	}

	return &RegistrationStatus{
		HasUsers:          !isFirstUser,
		AllowRegistration: s.allowRegistration,
	}, nil
}

func (s *jwtService) CreateInitialAdmin(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return nil
	}

	normalizedUsername, err := normalizeUsername(username)
	if err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     normalizedUsername,
		PasswordHash: hashedPassword,
	}

	if _, err := s.userRepo.CreateInitialAdminIfNoUsers(ctx, user); err != nil {
		return fmt.Errorf("failed to create initial admin user: %w", err)
	}
	return nil
}

func (s *jwtService) SetAuthCookies(w http.ResponseWriter, tokens *Tokens) {
	sameSite := http.SameSiteLaxMode
	if s.secureCookies {
		sameSite = http.SameSiteStrictMode
	}

	// #nosec G124 -- Secure is configurable for local HTTP/test; production defaults it to true.
	http.SetCookie(w, &http.Cookie{
		Name:     accessTokenCookie,
		Value:    tokens.AccessToken,
		Path:     "/",
		Domain:   s.cookieDomain,
		MaxAge:   int(s.accessExpiry.Seconds()),
		HttpOnly: true,
		Secure:   s.secureCookies,
		SameSite: sameSite,
	})

	// #nosec G124 -- Secure is configurable for local HTTP/test; production defaults it to true.
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    tokens.RefreshToken,
		Path:     "/",
		Domain:   s.cookieDomain,
		MaxAge:   int(s.refreshExpiry.Seconds()),
		HttpOnly: true,
		Secure:   s.secureCookies,
		SameSite: sameSite,
	})
}

func (s *jwtService) ClearAuthCookies(w http.ResponseWriter) {
	sameSite := http.SameSiteLaxMode
	if s.secureCookies {
		sameSite = http.SameSiteStrictMode
	}

	for _, name := range []string{accessTokenCookie, refreshTokenCookie} {
		// #nosec G124 -- Secure is configurable for local HTTP/test; production defaults it to true.
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Domain:   s.cookieDomain,
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   s.secureCookies,
			SameSite: sameSite,
		})
	}
}

func (s *jwtService) getTokensFromRequest(r *http.Request) (accessToken, refreshToken string) {
	if cookie, err := r.Cookie(accessTokenCookie); err == nil {
		accessToken = cookie.Value
	}
	if cookie, err := r.Cookie(refreshTokenCookie); err == nil {
		refreshToken = cookie.Value
	}
	return
}

func (s *jwtService) isFirstUser(ctx context.Context) (bool, error) {
	users, err := s.userRepo.List(ctx, 1, 0)
	if err != nil {
		return false, fmt.Errorf("failed to list users: %w", err)
	}
	return len(users) == 0, nil
}

func (s *jwtService) generateTokens(user *User) (*Tokens, error) {
	accessToken, err := s.jwt.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func normalizeUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" || len([]byte(username)) > maxUsernameBytes {
		return "", ErrInvalidUsername
	}
	if strings.ContainsFunc(username, unicode.IsControl) {
		return "", ErrInvalidUsername
	}
	return username, nil
}

func validatePassword(password string) error {
	passwordBytes := len([]byte(password))
	if passwordBytes < minPasswordBytes || passwordBytes > maxPasswordBytes {
		return ErrInvalidPassword
	}
	if strings.ContainsRune(password, '\x00') {
		return ErrInvalidPassword
	}
	return nil
}
