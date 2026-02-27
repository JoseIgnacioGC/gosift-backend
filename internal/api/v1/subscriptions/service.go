package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/JoseIgnacioGC/gosift-backend/internal/api/v1/users"
	"github.com/JoseIgnacioGC/gosift-backend/internal/feed"
)

var (
	errFeedAlreadySubscribed = errors.New("already subscribed to this feed")
	errInvalidFeed           = errors.New("invalid or unreachable feed URL")
)

type service struct {
	subRepo *Repository
}

func newService(subRepo *Repository) *service {
	return &service{subRepo: subRepo}
}

func (s *service) create(ctx context.Context, userID string, req CreateRequestDto) (*ResponseDto, error) {
	userPK := users.PKPrefix + userID

	existing, err := s.subRepo.FindByUserAndFeedURL(ctx, userPK, req.FeedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing subscription: %w", err)
	}
	if existing != nil {
		return nil, errFeedAlreadySubscribed
	}

	feedInfo, err := feed.Fetch(ctx, req.FeedURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidFeed, err.Error())
	}

	title := req.Title
	if title == "" {
		title = feedInfo.Title
	}

	siteURL := req.SiteURL
	if siteURL == "" {
		siteURL = feedInfo.SiteURL
	}

	subID := uuid.New().String()
	now := time.Now().UTC()

	sub := &Subscription{
		PK:        userPK,
		SK:        SKPrefix + subID,
		FeedURL:   req.FeedURL,
		Title:     title,
		SiteURL:   siteURL,
		Category:  req.Category,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.subRepo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return toResponseDto(subID, sub), nil
}

func (s *service) list(ctx context.Context, userID string) ([]ResponseDto, error) {
	userPK := users.PKPrefix + userID

	subs, err := s.subRepo.ListByUser(ctx, userPK)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	result := make([]ResponseDto, 0, len(subs))
	for i := range subs {
		subID := strings.TrimPrefix(subs[i].SK, SKPrefix)
		result = append(result, *toResponseDto(subID, &subs[i]))
	}

	return result, nil
}

func (s *service) update(ctx context.Context, userID string, subID string, req UpdateRequestDto) error {
	userPK := users.PKPrefix + userID
	sk := SKPrefix + subID

	now := time.Now().UTC()
	fields := &SubscriptionUpdate{
		Title:     req.Title,
		Category:  req.Category,
		IsActive:  req.IsActive,
		UpdatedAt: &now,
	}

	if err := s.subRepo.UpdateFields(ctx, userPK, sk, fields); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

func (s *service) delete(ctx context.Context, userID string, subID string) error {
	userPK := users.PKPrefix + userID
	sk := SKPrefix + subID

	if err := s.subRepo.Delete(ctx, userPK, sk); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

func toResponseDto(subID string, sub *Subscription) *ResponseDto {
	return &ResponseDto{
		ID:        subID,
		FeedURL:   sub.FeedURL,
		Title:     sub.Title,
		SiteURL:   sub.SiteURL,
		Category:  sub.Category,
		IsActive:  sub.IsActive,
		CreatedAt: sub.CreatedAt.Format(time.RFC3339),
		UpdatedAt: sub.UpdatedAt.Format(time.RFC3339),
	}
}
