package main

import (
	"github.com/newrelic/newrelic-client-go/v2/pkg/usermanagement"
)

type authenticationDomainsResponse struct {
	Actor usermanagement.Actor `json:"actor"`
}

type ResourceUser struct {
	id                       string
	name                     string
	email_id                 string
	authentication_domain_id string
	user_type                string
}

type ResourceGroup struct {
	id                       string
	name                     string
	authentication_domain_id string
	user_ids                 []string
}

var userTier = map[string]string{
	"Basic":         "BASIC_USER_TIER",
	"Core":          "CORE_USER_TIER",
	"Full platform": "FULL_USER_TIER",
}
