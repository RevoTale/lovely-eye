package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
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
	minPasswordBytes = 8
	maxPasswordBytes = 72
	maxUsernameBytes = 128
)

type Config struct {
	JWTSecret         string
	AccessTokenExpiry time.Duration
	RefreshExpiry     time.Duration
	AllowRegistration bool
}

// Service implements authentication and token lifecycle policy.
type Service struct {
	userStore         UserStore
	jwt               *jwtProvider
	allowRegistration bool
}

// NewService creates a new authentication service.
func NewService(userStore UserStore, cfg Config) *Service {
	return &Service{
		userStore:         userStore,
		jwt:               newJWTProvider(cfg.JWTSecret, cfg.AccessTokenExpiry, cfg.RefreshExpiry),
		allowRegistration: cfg.AllowRegistration,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*User, *Tokens, error) {
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

	storedUser := &StoredUser{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	if err := s.userStore.CreateForRegistration(ctx, storedUser, s.allowRegistration); err != nil {
		if errors.Is(err, ErrUserExists) {
			return nil, nil, ErrUserExists
		}
		if errors.Is(err, ErrRegistrationDisabled) {
			return nil, nil, ErrRegistrationDisabled
		}
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	user := publicUser(storedUser)

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*User, *Tokens, error) {
	username, err := normalizeUsername(input.Username)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	if err := validatePassword(input.Password); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	storedUser, err := s.userStore.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	if !checkPassword(input.Password, storedUser.PasswordHash) {
		return nil, nil, ErrInvalidCredentials
	}

	user := publicUser(storedUser)

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*Tokens, error) {
	claims, err := s.jwt.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	storedUser, err := s.userStore.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return s.generateTokens(publicUser(storedUser))
}

func (s *Service) ValidateAccessToken(token string) (*Claims, error) {
	return s.jwt.validateAccessToken(token)
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (*User, error) {
	storedUser, err := s.userStore.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return publicUser(storedUser), nil
}

func (s *Service) RegistrationStatus(ctx context.Context) (*RegistrationStatus, error) {
	isFirstUser, err := s.isFirstUser(ctx)
	if err != nil {
		return nil, err
	}

	return &RegistrationStatus{
		HasUsers:          !isFirstUser,
		AllowRegistration: s.allowRegistration,
	}, nil
}

func (s *Service) CreateInitialAdmin(ctx context.Context, username, password string) error {
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

	user := &StoredUser{
		Username:     normalizedUsername,
		PasswordHash: hashedPassword,
	}

	if _, err := s.userStore.CreateInitialAdminIfNoUsers(ctx, user); err != nil {
		return fmt.Errorf("failed to create initial admin user: %w", err)
	}
	return nil
}

func (s *Service) isFirstUser(ctx context.Context) (bool, error) {
	hasUsers, err := s.userStore.HasUsers(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check users: %w", err)
	}
	return !hasUsers, nil
}

func (s *Service) generateTokens(user *User) (*Tokens, error) {
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

func publicUser(user *StoredUser) *User {
	return &User{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
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
