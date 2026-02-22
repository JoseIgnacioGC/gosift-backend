package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/JoseIgnacioGC/gosift-backend/internal/api/v1/users"
	gosiftjwt "github.com/JoseIgnacioGC/gosift-backend/internal/platform/jwt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
)

const (
	bcryptCost      = 12
	defaultTimezone = "UTC"
	defaultLanguage = "en"
)

type Service struct {
	userRepo  *users.Repository
	jwtSecret string
}

func NewService(userRepo *users.Repository, jwtSecret string) *Service {
	return &Service{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *Service) Register(ctx context.Context, req RegisterRequestDto) (*RegisterResponseDto, error) {
	existing, err := s.userRepo.FindByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	timezone := req.Timezone
	if timezone == "" {
		timezone = defaultTimezone
	}
	language := req.Language
	if language == "" {
		language = defaultLanguage
	}

	userID := uuid.New().String()
	now := time.Now().UTC()

	user := &users.User{
		PK:            users.PKPrefix + userID,
		SK:            users.SKMetadata,
		Email:         strings.ToLower(req.Email),
		Username:      req.Username,
		PasswordHash:  string(hashedPassword),
		Name:          req.Name,
		AvatarURL:     "",
		Timezone:      timezone,
		Language:      language,
		EmailVerified: false,
		LastLoginAt:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := gosiftjwt.GenerateToken(userID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &RegisterResponseDto{
		ID:        userID,
		Email:     user.Email,
		Username:  user.Username,
		Name:      user.Name,
		Timezone:  user.Timezone,
		Language:  user.Language,
		CreatedAt: now.Format(time.RFC3339),
		Token:     token,
	}, nil
}
