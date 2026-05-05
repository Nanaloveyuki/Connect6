package service

import (
	"context"
	"strings"

	"connect6/backend/internal/domain/graph"
	"connect6/backend/internal/repository"
)

type PathService struct {
	relationships repository.RelationshipRepository
	users         *UserService
	defaultDepth  int
}

func NewPathService(relationships repository.RelationshipRepository, users *UserService, defaultDepth int) *PathService {
	return &PathService{relationships: relationships, users: users, defaultDepth: defaultDepth}
}

func (s *PathService) FindShortestPath(ctx context.Context, source string, target string, maxDepth int) (graph.Path, error) {
	source = normalizeUsername(source)
	target = normalizeUsername(target)
	if source == "" || target == "" {
		return graph.Path{}, graph.ErrPathNotFound
	}

	if maxDepth <= 0 {
		maxDepth = s.defaultDepth
	}

	if source == target {
		node := s.buildNode(ctx, source)
		return graph.Path{Source: source, Target: target, Degree: 0, Meta: graph.SourceMeta{Mode: "direct", Reason: "same-user"}, Nodes: []graph.Node{node}, Edges: []graph.Edge{}, Steps: []graph.Step{}}, nil
	}

	queue := []string{source}
	visited := map[string]bool{source: true}
	depth := map[string]int{source: 0}
	previous := map[string]string{}
	relations := map[string]string{}
	pathMode := "seed"
	pathReason := "local-development"
	found := false

	for len(queue) > 0 && !found {
		current := queue[0]
		queue = queue[1:]

		if depth[current] >= maxDepth {
			continue
		}

		neighbors, mode, reason, err := s.relationships.GetNeighbors(ctx, current)
		if err != nil {
			return graph.Path{}, err
		}

		if depth[current] == 0 {
			pathMode = mode
			pathReason = reason
		}

		for _, neighbor := range neighbors {
			normalizedNeighbor := normalizeUsername(neighbor.Username)
			if visited[normalizedNeighbor] {
				continue
			}

			visited[normalizedNeighbor] = true
			previous[normalizedNeighbor] = current
			relations[normalizedNeighbor] = neighbor.Relation
			depth[normalizedNeighbor] = depth[current] + 1
			queue = append(queue, normalizedNeighbor)

			if normalizedNeighbor == target {
				found = true
				break
			}
		}
	}

	if !found {
		return graph.Path{}, graph.ErrPathNotFound
	}

	usernames := buildUsernamePath(previous, source, target)
	nodes := make([]graph.Node, 0, len(usernames))
	edges := make([]graph.Edge, 0, len(usernames)-1)
	steps := make([]graph.Step, 0, len(usernames)-1)

	for _, username := range usernames {
		nodes = append(nodes, s.buildNode(ctx, username))
	}

	for index := 0; index < len(usernames)-1; index++ {
		relation := relations[usernames[index+1]]
		if relation == "" {
			relation = "related"
		}
		edges = append(edges, graph.Edge{Source: usernames[index], Target: usernames[index+1], Label: relation})
		steps = append(steps, graph.Step{From: usernames[index], To: usernames[index+1], Relation: relation})
	}

	return graph.Path{Source: source, Target: target, Degree: len(usernames) - 1, Meta: graph.SourceMeta{Mode: pathMode, Reason: pathReason}, Nodes: nodes, Edges: edges, Steps: steps}, nil
}

func (s *PathService) buildNode(ctx context.Context, username string) graph.Node {
	user, err := s.users.GetUser(ctx, username)
	if err != nil {
		return graph.Node{ID: username, Label: username, ProfileURL: "https://github.com/" + username}
	}

	label := user.Login
	if strings.TrimSpace(user.Name) != "" {
		label = user.Name
	}

	return graph.Node{ID: user.Login, Label: label, AvatarURL: user.AvatarURL, ProfileURL: user.ProfileURL}
}

func buildUsernamePath(previous map[string]string, source string, target string) []string {
	path := []string{target}
	for current := target; current != source; {
		current = previous[current]
		path = append(path, current)
	}

	for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
		path[left], path[right] = path[right], path[left]
	}

	return path
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}
