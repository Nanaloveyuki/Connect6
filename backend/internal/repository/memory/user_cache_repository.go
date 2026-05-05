package memory

import (
	"context"
	"strings"
	"sync"

	domaingithub "connect6/backend/internal/domain/github"
)

type UserCacheRepository struct {
	mu    sync.RWMutex
	users map[string]domaingithub.User
}

func NewUserCacheRepository() *UserCacheRepository {
	return &UserCacheRepository{users: make(map[string]domaingithub.User)}
}

func (r *UserCacheRepository) Get(_ context.Context, username string) (domaingithub.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[cacheKey(username)]
	return user, ok
}

func (r *UserCacheRepository) Set(_ context.Context, user domaingithub.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[cacheKey(user.Login)] = user
	return nil
}

func cacheKey(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}
