package service

import (
	"context"
	"testing"

	domaingithub "connect6/backend/internal/domain/github"
	"connect6/backend/internal/repository/memory"
)

type stubUserProvider struct{}

func (stubUserProvider) FetchUser(_ context.Context, username string) (domaingithub.User, error) {
	return domaingithub.User{
		Login:      username,
		Name:       username,
		ProfileURL: "https://github.com/" + username,
	}, nil
}

func TestGetUserCachesProviderResult(t *testing.T) {
	t.Parallel()

	cache := memory.NewUserCacheRepository()
	service := NewUserService(cache, stubUserProvider{})

	user, err := service.GetUser(context.Background(), "torvalds")
	if err != nil {
		t.Fatalf("GetUser returned error: %v", err)
	}

	if user.Login != "torvalds" {
		t.Fatalf("expected login torvalds, got %s", user.Login)
	}

	cachedUser, ok := cache.Get(context.Background(), "torvalds")
	if !ok {
		t.Fatal("expected user to be cached")
	}

	if cachedUser.Login != "torvalds" {
		t.Fatalf("expected cached login torvalds, got %s", cachedUser.Login)
	}
}
