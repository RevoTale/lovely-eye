package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserExists           = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrRegistrationDisabled = errors.New("registration is disabled")
)

// Cookie and header names
const (
	accessTokenCookie  = "le_access"
	refreshTokenCookie = "le_refresh"
	csrfTokenCookie    = "le_csrf"
	csrfHeaderName     = "X-CSRF-Token"
)

// Config contains authentication configuration.
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
	isFirstUser, err := s.isFirstUser(ctx)
	if err != nil {
		return nil, nil, err
	}

	if !isFirstUser && !s.allowRegistration {
		return nil, nil, ErrRegistrationDisabled
	}

	existing, _ := s.userRepo.GetByUsername(ctx, input.Username)
	if existing != nil {
		return nil, nil, ErrUserExists
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return nil, nil, err
	}

	role := "user"
	if isFirstUser {
		role = "admin"
	}

	dbUser := &models.User{
		Username:     input.Username,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	if err := s.userRepo.Create(ctx, dbUser); err != nil {
		return nil, nil, err
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
	dbUser, err := s.userRepo.GetByUsername(ctx, input.Username)
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

func (s *jwtService) CreateInitialAdmin(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return nil
	}

	isFirst, err := s.isFirstUser(ctx)
	if err != nil {
		return err
	}
	if !isFirst {
		return nil
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         "admin",
	}

	return s.userRepo.Create(ctx, user)
}

func (s *jwtService) SetAuthCookies(w http.ResponseWriter, tokens *Tokens) {
	// Use Lax for development (allows cookies across localhost ports)
	// Use Strict for production (same origin only)
	sameSite := http.SameSiteLaxMode
	if s.secureCookies {
		sameSite = http.SameSiteStrictMode
	}

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

	http.SetCookie(w, &http.Cookie{
		Name:     csrfTokenCookie,
		Value:    tokens.CSRFToken,
		Path:     "/",
		Domain:   s.cookieDomain,
		MaxAge:   int(s.refreshExpiry.Seconds()),
		HttpOnly: false,
		Secure:   s.secureCookies,
		SameSite: sameSite,
	})
}

func (s *jwtService) ClearAuthCookies(w http.ResponseWriter) {
	sameSite := http.SameSiteLaxMode
	if s.secureCookies {
		sameSite = http.SameSiteStrictMode
	}

	for _, name := range []string{accessTokenCookie, refreshTokenCookie, csrfTokenCookie} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Domain:   s.cookieDomain,
			MaxAge:   -1,
			HttpOnly: name != csrfTokenCookie,
			Secure:   s.secureCookies,
			SameSite: sameSite,
		})
	}
}

// getTokensFromRequest extracts tokens from cookies.
func (s *jwtService) getTokensFromRequest(r *http.Request) (accessToken, refreshToken string) {
	if cookie, err := r.Cookie(accessTokenCookie); err == nil {
		accessToken = cookie.Value
	}
	if cookie, err := r.Cookie(refreshTokenCookie); err == nil {
		refreshToken = cookie.Value
	}
	return
}

// validateCSRF validates the CSRF token using constant-time comparison.
func (s *jwtService) validateCSRF(r *http.Request) bool {
	cookie, err := r.Cookie(csrfTokenCookie)
	if err != nil {
		return false
	}

	headerToken := r.Header.Get(csrfHeaderName)
	if headerToken == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(headerToken)) == 1
}

func (s *jwtService) generateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return base64.URLEncoding.EncodeToString([]byte(time.Now().String()))
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (s *jwtService) isFirstUser(ctx context.Context) (bool, error) {
	users, err := s.userRepo.List(ctx, 1, 0)
	if err != nil {
		return false, err
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
		CSRFToken:    s.generateCSRFToken(),
	}, nil
}
