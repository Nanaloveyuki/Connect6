package githubapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"connect6/backend/internal/config"
	"connect6/backend/internal/constant"
	domaingithub "connect6/backend/internal/domain/github"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	token          string
	resty          *resty.Client
	followingLimit int
	orgLimit       int
}

type userResponse struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Bio       string `json:"bio"`
}

type relationshipUserResponse struct {
	Login string `json:"login"`
}

type organizationResponse struct {
	Login string `json:"login"`
}

func NewClient(cfg config.GitHubConfig) *Client {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}

	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("User-Agent", "Connect6").
		SetTimeout(time.Duration(cfg.RequestTimeoutSeconds) * time.Second)

	return &Client{
		token:          cfg.Token,
		resty:          client,
		followingLimit: cfg.FollowingLimit,
		orgLimit:       cfg.OrgLimit,
	}
}

func (c *Client) FetchUser(ctx context.Context, username string) (domaingithub.User, error) {
	var payload userResponse
	request := c.resty.R().SetContext(ctx).SetResult(&payload)
	if strings.TrimSpace(c.token) != "" {
		request.SetAuthToken(c.token)
	}

	response, err := request.Get("/users/" + username)
	if err != nil {
		return domaingithub.User{}, err
	}

	switch response.StatusCode() {
	case http.StatusOK:
		return domaingithub.User{Login: payload.Login, Name: payload.Name, AvatarURL: payload.AvatarURL, ProfileURL: payload.HTMLURL, Bio: payload.Bio}, nil
	case http.StatusNotFound:
		return domaingithub.User{}, domaingithub.ErrUserNotFound
	default:
		return domaingithub.User{}, fmt.Errorf("github api returned status %d", response.StatusCode())
	}
}

func (c *Client) FetchRelationships(ctx context.Context, username string) ([]domaingithub.Relationship, error) {
	following, err := c.fetchFollowing(ctx, username)
	if err != nil {
		return nil, err
	}

	_, err = c.fetchOrganizations(ctx, username)
	if err != nil {
		return nil, err
	}

	relationships := make([]domaingithub.Relationship, 0, len(following))
	for _, login := range following {
		relationships = append(relationships, domaingithub.Relationship{Username: login, Relation: constant.RelationFollows})
	}

	return relationships, nil
}

func (c *Client) fetchFollowing(ctx context.Context, username string) ([]string, error) {
	var payload []relationshipUserResponse
	request := c.resty.R().SetContext(ctx).SetResult(&payload).SetQueryParam("per_page", fmt.Sprintf("%d", c.followingLimit))
	if strings.TrimSpace(c.token) != "" {
		request.SetAuthToken(c.token)
	}

	response, err := request.Get("/users/" + username + "/following")
	if err != nil {
		return nil, err
	}

	switch response.StatusCode() {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, domaingithub.ErrUserNotFound
	case http.StatusUnauthorized:
		return nil, domaingithub.ErrRelationshipProviderUnauthorized
	case http.StatusForbidden, http.StatusTooManyRequests:
		return nil, domaingithub.ErrRelationshipProviderRateLimited
	default:
		return nil, fmt.Errorf("github following api returned status %d", response.StatusCode())
	}

	result := make([]string, 0, len(payload))
	for _, item := range payload {
		if trimmed := strings.TrimSpace(item.Login); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result, nil
}

func (c *Client) fetchOrganizations(ctx context.Context, username string) ([]string, error) {
	var payload []organizationResponse
	request := c.resty.R().SetContext(ctx).SetResult(&payload).SetQueryParam("per_page", fmt.Sprintf("%d", c.orgLimit))
	if strings.TrimSpace(c.token) != "" {
		request.SetAuthToken(c.token)
	}

	response, err := request.Get("/users/" + username + "/orgs")
	if err != nil {
		return nil, err
	}

	switch response.StatusCode() {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, domaingithub.ErrUserNotFound
	case http.StatusUnauthorized:
		return nil, domaingithub.ErrRelationshipProviderUnauthorized
	case http.StatusForbidden, http.StatusTooManyRequests:
		return nil, domaingithub.ErrRelationshipProviderRateLimited
	default:
		return nil, fmt.Errorf("github orgs api returned status %d", response.StatusCode())
	}

	result := make([]string, 0, len(payload))
	for _, item := range payload {
		if trimmed := strings.TrimSpace(item.Login); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result, nil
}
