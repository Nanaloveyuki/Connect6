package repository

import (
	"context"

	domaingithub "connect6/backend/internal/domain/github"
)

type UserProvider interface {
	FetchUser(ctx context.Context, username string) (domaingithub.User, error)
}

type RelationshipProvider interface {
	FetchRelationships(ctx context.Context, username string) ([]domaingithub.Relationship, error)
}

type UserCacheRepository interface {
	Get(ctx context.Context, username string) (domaingithub.User, bool)
	Set(ctx context.Context, user domaingithub.User) error
}

type RelationshipRepository interface {
	GetNeighbors(ctx context.Context, username string) ([]domaingithub.Relationship, string, string, error)
}
