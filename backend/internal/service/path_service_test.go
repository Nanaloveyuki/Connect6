package service

import (
	"context"
	"testing"

	"connect6/backend/internal/domain/graph"
	"connect6/backend/internal/repository/memory"
)

type stubUserService struct{}

func TestFindShortestPathReturnsSeedPath(t *testing.T) {
	t.Parallel()

	relationships := memory.NewRelationshipRepository()
	users := NewUserService(memory.NewUserCacheRepository(), stubUserProvider{})
	service := NewPathService(relationships, users, 6)

	path, err := service.FindShortestPath(context.Background(), "torvalds", "antfu", 6)
	if err != nil {
		t.Fatalf("FindShortestPath returned error: %v", err)
	}

	if path.Degree != 2 {
		t.Fatalf("expected degree 2, got %d", path.Degree)
	}

	if len(path.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(path.Nodes))
	}

	if path.Source != "torvalds" || path.Target != "antfu" {
		t.Fatalf("unexpected endpoints: %+v", path)
	}

	if path.Meta.Mode == "" {
		t.Fatal("expected path source metadata to be populated")
	}
}

func TestFindShortestPathReturnsNotFound(t *testing.T) {
	t.Parallel()

	relationships := memory.NewRelationshipRepository()
	users := NewUserService(memory.NewUserCacheRepository(), stubUserProvider{})
	service := NewPathService(relationships, users, 1)

	_, err := service.FindShortestPath(context.Background(), "torvalds", "antfu", 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != graph.ErrPathNotFound {
		t.Fatalf("expected ErrPathNotFound, got %v", err)
	}
}

func TestFindShortestPathReturnsDirectMetaForSameUser(t *testing.T) {
	t.Parallel()

	relationships := memory.NewRelationshipRepository()
	users := NewUserService(memory.NewUserCacheRepository(), stubUserProvider{})
	service := NewPathService(relationships, users, 6)

	path, err := service.FindShortestPath(context.Background(), "torvalds", "torvalds", 6)
	if err != nil {
		t.Fatalf("FindShortestPath returned error: %v", err)
	}

	if path.Meta.Mode != "direct" || path.Meta.Reason != "same-user" {
		t.Fatalf("unexpected direct path metadata: %+v", path.Meta)
	}
}
