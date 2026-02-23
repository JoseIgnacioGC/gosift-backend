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
	errEmailAlreadyExists = errors.New("email already registered")
)

const (
	bcryptCost      = 12
	defaultTimezone = "UTC"
	defaultLanguage = "en"
)

type service struct {
	userRepo  *users.Repository
	jwtSecret string
}

func newService(userRepo *users.Repository, jwtSecret string) *service {
	return &service{userRepo: userRepo, jwtSecret: jwtSecret}
}

// TODO: Implement email verification to prevent timing attacks. Send a verification link to the user's email.
func (s *service) register(ctx context.Context, req RegisterRequestDto) (*ResponseDto, error) {
	existing, err := s.userRepo.FindByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existing != nil {
		return nil, errEmailAlreadyExists
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

	return &ResponseDto{
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

const DUMMY_PASSWORD_HASH = "$2a$12$wJjH6Zy1G8Q9s5nXl3e7uJ8v1b9c0dE4f5g6h7i8j9k0l1m2n3o4p5q"

func (s *service) login(ctx context.Context, req LoginRequestDto) (*ResponseDto, error) {
	existing, err := s.userRepo.FindByEmail(ctx, strings.ToLower(req.Email))
	if err != nil || existing == nil {
		_ = bcrypt.CompareHashAndPassword([]byte(DUMMY_PASSWORD_HASH), []byte(req.Password))

		return nil, fmt.Errorf("Invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("Invalid email or password")
	}

	now := time.Now().UTC()
	err = s.userRepo.UpdateFields(ctx, existing.PK, &users.UserUpdate{
		LastLoginAt: &now,
		UpdatedAt:   &now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update login timestamp: %w", err)
	}

	userID := strings.TrimPrefix(existing.PK, users.PKPrefix)

	token, err := gosiftjwt.GenerateToken(userID, existing.Email, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &ResponseDto{
		ID:        userID,
		Email:     existing.Email,
		Username:  existing.Username,
		Name:      existing.Name,
		Timezone:  existing.Timezone,
		Language:  existing.Language,
		CreatedAt: existing.CreatedAt.Format(time.RFC3339),
		Token:     token,
	}, nil
}
