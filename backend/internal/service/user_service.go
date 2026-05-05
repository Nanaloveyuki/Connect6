package service

import (
	"context"
	"strings"

	domaingithub "connect6/backend/internal/domain/github"
	"connect6/backend/internal/repository"
)

type UserService struct {
	cache    repository.UserCacheRepository
	provider repository.UserProvider
}

func NewUserService(cache repository.UserCacheRepository, provider repository.UserProvider) *UserService {
	return &UserService{cache: cache, provider: provider}
}

func (s *UserService) GetUser(ctx context.Context, username string) (domaingithub.User, error) {
	normalizedUsername := strings.ToLower(strings.TrimSpace(username))
	if normalizedUsername == "" {
		return domaingithub.User{}, domaingithub.ErrUserNotFound
	}

	if user, ok := s.cache.Get(ctx, normalizedUsername); ok {
		return user, nil
	}

	user, err := s.provider.FetchUser(ctx, normalizedUsername)
	if err != nil {
		return domaingithub.User{}, err
	}

	if err := s.cache.Set(ctx, user); err != nil {
		return domaingithub.User{}, err
	}

	return user, nil
}

