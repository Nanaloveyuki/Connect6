package memory

import (
	"context"
	"strings"

	"connect6/backend/internal/constant"
	domaingithub "connect6/backend/internal/domain/github"
)

type RelationshipRepository struct {
	adjacency map[string][]domaingithub.Relationship
}

func NewRelationshipRepository() *RelationshipRepository {
	return &RelationshipRepository{
		adjacency: map[string][]domaingithub.Relationship{
			"torvalds": {
				{Username: "gaearon", Relation: constant.RelationMaintainerBridge},
				{Username: "yyx990803", Relation: constant.RelationFrameworkCircle},
			},
			"gaearon": {
				{Username: "torvalds", Relation: constant.RelationMaintainerBridge},
				{Username: "sindresorhus", Relation: constant.RelationPackageLine},
				{Username: "tj", Relation: constant.RelationMaintainerBridge},
			},
			"sindresorhus": {
				{Username: "gaearon", Relation: constant.RelationPackageLine},
				{Username: "antfu", Relation: constant.RelationPackageLine},
			},
			"antfu": {
				{Username: "sindresorhus", Relation: constant.RelationPackageLine},
				{Username: "yyx990803", Relation: constant.RelationFrameworkCircle},
			},
			"yyx990803": {
				{Username: "torvalds", Relation: constant.RelationFrameworkCircle},
				{Username: "antfu", Relation: constant.RelationFrameworkCircle},
			},
			"tj": {
				{Username: "gaearon", Relation: constant.RelationMaintainerBridge},
			},
		},
	}
}

func (r *RelationshipRepository) GetNeighbors(_ context.Context, username string) ([]domaingithub.Relationship, string, string, error) {
	key := normalizeUsername(username)
	neighbors := r.adjacency[key]
	result := make([]domaingithub.Relationship, len(neighbors))
	copy(result, neighbors)
	return result, "seed", "local-development", nil
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}
