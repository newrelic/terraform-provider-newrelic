package http

import (
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// RequestAuthorizer is an interface that allows customizatino of how a request is authorized.
type RequestAuthorizer interface {
	AuthorizeRequest(*retryablehttp.Request, *config.Config)
}

// NerdGraphAuthorizer authorizes calls to NerdGraph.
type NerdGraphAuthorizer struct{}

// AuthorizeRequest is responsible for setting up auth for a request.
func (a *NerdGraphAuthorizer) AuthorizeRequest(req *retryablehttp.Request, c *config.Config) {
	req.Header.Set("Api-Key", c.PersonalAPIKey)
}

// PersonalAPIKeyCapableV2Authorizer authorizes V2 endpoints that can use a personal API key.
type PersonalAPIKeyCapableV2Authorizer struct{}

// AuthorizeRequest is responsible for setting up auth for a request.
func (a *PersonalAPIKeyCapableV2Authorizer) AuthorizeRequest(req *retryablehttp.Request, c *config.Config) {
	if c.PersonalAPIKey != "" {
		req.Header.Set("Api-Key", c.PersonalAPIKey)
		req.Header.Set("Auth-Type", "User-Api-Key")
	} else {
		req.Header.Set("X-Api-Key", c.AdminAPIKey)
	}
}

// ClassicV2Authorizer authorizes V2 endpoints that cannot use a personal API key.
type ClassicV2Authorizer struct{}

// AuthorizeRequest is responsible for setting up auth for a request.
func (a *ClassicV2Authorizer) AuthorizeRequest(req *retryablehttp.Request, c *config.Config) {
	req.Header.Set("X-Api-Key", c.AdminAPIKey)
}
