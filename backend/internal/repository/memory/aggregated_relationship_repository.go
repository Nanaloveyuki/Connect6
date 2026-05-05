package memory

import (
	"context"
	"errors"

	domaingithub "connect6/backend/internal/domain/github"
	"connect6/backend/internal/repository"
)

type AggregatedRelationshipRepository struct {
	seed     repository.RelationshipRepository
	provider repository.RelationshipProvider
}

func NewAggregatedRelationshipRepository(seed repository.RelationshipRepository, provider repository.RelationshipProvider) *AggregatedRelationshipRepository {
	return &AggregatedRelationshipRepository{seed: seed, provider: provider}
}

func (r *AggregatedRelationshipRepository) GetNeighbors(ctx context.Context, username string) ([]domaingithub.Relationship, string, string, error) {
	providerRelationships, err := r.provider.FetchRelationships(ctx, username)
	if err == nil {
		return providerRelationships, "live", "github-following", nil
	}

	reason := ""
	switch {
	case errors.Is(err, domaingithub.ErrRelationshipProviderUnauthorized):
		reason = "github-unauthorized"
	case errors.Is(err, domaingithub.ErrRelationshipProviderRateLimited):
		reason = "github-rate-limited"
	case errors.Is(err, domaingithub.ErrRelationshipProviderUnavailable):
		reason = "github-provider-unavailable"
	}

	if reason == "" {
		return nil, "", "", err
	}

	seedRelationships, mode, _, seedErr := r.seed.GetNeighbors(ctx, username)
	if seedErr != nil {
		return nil, "", "", seedErr
	}

	return seedRelationships, mode, reason, nil
}
