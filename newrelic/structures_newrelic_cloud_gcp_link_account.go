package newrelic

import "encoding/json"

// Auxiliary helpers for the newrelic_cloud_gcp_link_account resource (WIF credential
// construction). Kept out of the resource file so that file holds ~CRUD only.

// gcpWIFOIDCEndpoint returns the New Relic OIDC token endpoint for the given provider region.
// This URL is set as credential_source.url in the WIF credential JSON and tells GCP STS
// where to fetch the subject token from.
func gcpWIFOIDCEndpoint(region string) string {
	switch region {
	case "EU":
		return "https://oidc.eu.newrelic.com/r/gcp-cmp"
	case "JP":
		return "https://oidc.jp.newrelic.com/r/gcp-cmp"
	case "Staging":
		return "https://oidc-staging.newrelic.com/r/gcp-cmp"
	default: // US
		return "https://oidc.newrelic.com/r/gcp-cmp"
	}
}

// gcpBuildWIFCredential constructs the GCP Workload Identity Federation credential
// JSON string that cloudAuthenticateIntegration expects as its payload. All fixed
// fields (universe_domain, type, subject_token_type, token_url, format) are set to
// their required values; the caller supplies only the environment-specific inputs.
func gcpBuildWIFCredential(audience, serviceAccountEmail, region string) (string, error) {
	cred := map[string]interface{}{
		"universe_domain":    "googleapis.com",
		"type":               "external_account",
		"audience":           audience,
		"subject_token_type": "urn:ietf:params:oauth:token-type:jwt",
		"token_url":          "https://sts.googleapis.com/v1/token",
		"credential_source": map[string]interface{}{
			"url":     gcpWIFOIDCEndpoint(region),
			"headers": map[string]interface{}{},
			"format": map[string]interface{}{
				"type":                     "json",
				"subject_token_field_name": "access_token",
			},
		},
		"service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/" +
			serviceAccountEmail + ":generateAccessToken",
	}
	b, err := json.Marshal(cred)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
