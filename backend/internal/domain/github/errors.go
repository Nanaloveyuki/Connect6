package github

import "errors"

var ErrUserNotFound = errors.New("github user not found")
var ErrRelationshipProviderUnavailable = errors.New("github relationship provider unavailable")
var ErrRelationshipProviderUnauthorized = errors.New("github relationship provider unauthorized")
var ErrRelationshipProviderRateLimited = errors.New("github relationship provider rate limited")
