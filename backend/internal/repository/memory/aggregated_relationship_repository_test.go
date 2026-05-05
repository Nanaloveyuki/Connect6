package memory

import (
	"context"
	"errors"
	"testing"

	domaingithub "connect6/backend/internal/domain/github"
)

type stubRelationshipProvider struct {
	relationships []domaingithub.Relationship
	err           error
}

func (s stubRelationshipProvider) FetchRelationships(_ context.Context, _ string) ([]domaingithub.Relationship, error) {
	if s.err != nil {
		return nil, s.err
	}

	result := make([]domaingithub.Relationship, len(s.relationships))
	copy(result, s.relationships)
	return result, nil
}

func TestAggregatedRelationshipRepositoryPrefersProviderRelationships(t *testing.T) {
	t.Parallel()

	seed := NewRelationshipRepository()
	provider := stubRelationshipProvider{
		relationships: []domaingithub.Relationship{{Username: "alice", Relation: "follows"}},
	}
	repository := NewAggregatedRelationshipRepository(seed, provider)

	neighbors, mode, reason, err := repository.GetNeighbors(context.Background(), "torvalds")
	if err != nil {
		t.Fatalf("GetNeighbors returned error: %v", err)
	}

	if len(neighbors) != 1 || neighbors[0].Username != "alice" {
		t.Fatalf("expected provider relationships, got %+v", neighbors)
	}

	if mode != "live" || reason != "github-following" {
		t.Fatalf("unexpected provider source metadata: %s %s", mode, reason)
	}
}

func TestAggregatedRelationshipRepositoryFallsBackToSeedOnProviderError(t *testing.T) {
	t.Parallel()

	seed := NewRelationshipRepository()
	provider := stubRelationshipProvider{err: domaingithub.ErrRelationshipProviderUnavailable}
	repository := NewAggregatedRelationshipRepository(seed, provider)

	neighbors, mode, reason, err := repository.GetNeighbors(context.Background(), "torvalds")
	if err != nil {
		t.Fatalf("GetNeighbors returned error: %v", err)
	}

	if len(neighbors) == 0 {
		t.Fatal("expected seed fallback relationships")
	}

	if mode != "seed" || reason != "github-provider-unavailable" {
		t.Fatalf("unexpected fallback source metadata: %s %s", mode, reason)
	}
}

func TestAggregatedRelationshipRepositoryMapsUnauthorizedFallbackReason(t *testing.T) {
	t.Parallel()

	seed := NewRelationshipRepository()
	provider := stubRelationshipProvider{err: domaingithub.ErrRelationshipProviderUnauthorized}
	repository := NewAggregatedRelationshipRepository(seed, provider)

	_, mode, reason, err := repository.GetNeighbors(context.Background(), "torvalds")
	if err != nil {
		t.Fatalf("GetNeighbors returned error: %v", err)
	}

	if mode != "seed" || reason != "github-unauthorized" {
		t.Fatalf("unexpected unauthorized fallback metadata: %s %s", mode, reason)
	}
}

func TestAggregatedRelationshipRepositoryMapsRateLimitedFallbackReason(t *testing.T) {
	t.Parallel()

	seed := NewRelationshipRepository()
	provider := stubRelationshipProvider{err: domaingithub.ErrRelationshipProviderRateLimited}
	repository := NewAggregatedRelationshipRepository(seed, provider)

	_, mode, reason, err := repository.GetNeighbors(context.Background(), "torvalds")
	if err != nil {
		t.Fatalf("GetNeighbors returned error: %v", err)
	}

	if mode != "seed" || reason != "github-rate-limited" {
		t.Fatalf("unexpected rate-limit fallback metadata: %s %s", mode, reason)
	}
}

func TestAggregatedRelationshipRepositoryKeepsProviderEmptyResult(t *testing.T) {
	t.Parallel()

	seed := NewRelationshipRepository()
	provider := stubRelationshipProvider{relationships: []domaingithub.Relationship{}}
	repository := NewAggregatedRelationshipRepository(seed, provider)

	neighbors, mode, reason, err := repository.GetNeighbors(context.Background(), "torvalds")
	if err != nil {
		t.Fatalf("GetNeighbors returned error: %v", err)
	}

	if len(neighbors) != 0 {
		t.Fatalf("expected empty provider result to be preserved, got %+v", neighbors)
	}

	if mode != "live" || reason != "github-following" {
		t.Fatalf("unexpected empty-result source metadata: %s %s", mode, reason)
	}
}

func TestAggregatedRelationshipRepositoryReturnsUnexpectedProviderError(t *testing.T) {
	t.Parallel()

	seed := NewRelationshipRepository()
	provider := stubRelationshipProvider{err: errors.New("bad gateway")}
	repository := NewAggregatedRelationshipRepository(seed, provider)

	_, _, _, err := repository.GetNeighbors(context.Background(), "torvalds")
	if err == nil {
		t.Fatal("expected provider error to be returned")
	}
}
