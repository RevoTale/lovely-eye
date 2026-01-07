package services

import (
	"context"
	"errors"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
)

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserExists            = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrRegistrationDisabled  = errors.New("registration is disabled")
)

type AuthService struct {
	userRepo          *repository.UserRepository
	jwtService        *auth.JWTService
	allowRegistration bool
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *auth.JWTService, allowRegistration bool) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		jwtService:        jwtService,
		allowRegistration: allowRegistration,
	}
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type RegisterInput struct {
	Username string
	Password string
}

type LoginInput struct {
	Username string
	Password string
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*models.User, *AuthTokens, error) {
	// Check if this is the first user (will become admin)
	isFirstUser, err := s.isFirstUser(ctx)
	if err != nil {
		return nil, nil, err
	}

	// If not first user and registration is disabled, reject
	if !isFirstUser && !s.allowRegistration {
		return nil, nil, ErrRegistrationDisabled
	}

	// Check if user exists
	existing, _ := s.userRepo.GetByUsername(ctx, input.Username)
	if existing != nil {
		return nil, nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, nil, err
	}

	// First user becomes admin, others are regular users
	role := "user"
	if isFirstUser {
		role = "admin"
	}

	user := &models.User{
		Username:     input.Username,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*models.User, *AuthTokens, error) {
	user, err := s.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if !auth.CheckPassword(input.Password, user.PasswordHash) {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.generateTokens(user)
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*auth.Claims, error) {
	return s.jwtService.ValidateAccessToken(tokenString)
}

func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *AuthService) isFirstUser(ctx context.Context) (bool, error) {
	users, err := s.userRepo.List(ctx, 1, 0)
	if err != nil {
		return false, err
	}
	return len(users) == 0, nil
}

func (s *AuthService) generateTokens(user *models.User) (*AuthTokens, error) {
	accessToken, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// CreateInitialAdmin creates the initial admin user if no users exist
func (s *AuthService) CreateInitialAdmin(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return nil // No initial admin configured
	}

	isFirst, err := s.isFirstUser(ctx)
	if err != nil {
		return err
	}
	if !isFirst {
		return nil // Users already exist
	}

	hashedPassword, err := auth.HashPassword(password)
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
