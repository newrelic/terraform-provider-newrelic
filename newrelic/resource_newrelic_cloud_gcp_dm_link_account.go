package newrelic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

// gcpDmAuthenticateMutation calls cloudAuthenticateIntegration to obtain an authReferenceId
// for WIF-based GCP v2 account linking (Step 1 of 2-step Create).
const gcpDmAuthenticateMutation = `mutation(
	$accountId: Int!,
	$providerSlug: CloudProviderType!,
	$authType: AuthenticationType!,
	$payload: String!,
) {
	cloudAuthenticateIntegration(
		accountId: $accountId
		providerSlug: $providerSlug
		authType: $authType
		payload: $payload
	) {
		authReferenceId
	}
}`

// gcpDmLinkAccountMutation links a GCP project to a New Relic account using authReferenceId.
// We define this locally because cloud.CloudGcpLinkAccountInput does not include authReferenceId.
const gcpDmLinkAccountMutation = `mutation(
	$accountId: Int!,
	$accounts: CloudLinkCloudAccountsInput!,
) {
	cloudLinkAccount(
		accountId: $accountId
		accounts: $accounts
	) {
		linkedAccounts {
			id
			nrAccountId
			name
		}
		errors {
			type
			message
		}
	}
}`

// gcpDmGetLinkedAccountQuery fetches only the basic fields of a linked account
// (id, name, nrAccountId, externalId) without requesting integrations.
// This avoids the "Abstract type 'Integration' must resolve to an Object type"
// error that occurs when GetLinkedAccountWithContext encounters GCP v2-specific
// integration types that its inline fragments don't cover.
const gcpDmGetLinkedAccountQuery = `query($accountId: Int!, $linkedAccountId: Int!) {
	actor {
		account(id: $accountId) {
			cloud {
				linkedAccount(id: $linkedAccountId) {
					id
					name
					nrAccountId
					externalId
				}
			}
		}
	}
}`

// gcpDmLinkedAccountResp is the response type for gcpDmGetLinkedAccountQuery.
type gcpDmLinkedAccountResp struct {
	Actor struct {
		Account struct {
			Cloud struct {
				LinkedAccount *struct {
					ID          int    `json:"id"`
					Name        string `json:"name"`
					NrAccountId int    `json:"nrAccountId"`
					ExternalId  string `json:"externalId"`
				} `json:"linkedAccount"`
			} `json:"cloud"`
		} `json:"account"`
	} `json:"actor"`
}

// gcpDmAuthResp is the NerdGraph response for cloudAuthenticateIntegration.
type gcpDmAuthResp struct {
	CloudAuthenticateIntegration struct {
		AuthReferenceId string `json:"authReferenceId"`
	} `json:"cloudAuthenticateIntegration"`
}

// gcpDmLinkAccountInput is a local GCP link account input that includes authReferenceId,
// which is required for GCP v2 / WIF authentication but absent from cloud.CloudGcpLinkAccountInput.
type gcpDmLinkAccountInput struct {
	Name            string `json:"name"`
	ProjectId       string `json:"projectId"`
	AuthReferenceId string `json:"authReferenceId,omitempty"`
}

// gcpDmLinkResp is the NerdGraph response for cloudLinkAccount.
type gcpDmLinkResp struct {
	CloudLinkAccount struct {
		LinkedAccounts []struct {
			ID          int    `json:"id"`
			NrAccountId int    `json:"nrAccountId"`
			Name        string `json:"name"`
		} `json:"linkedAccounts"`
		Errors []struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"errors"`
	} `json:"cloudLinkAccount"`
}

func resourceNewRelicCloudGcpDmLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudGcpDmLinkAccountCreate,
		ReadContext:   resourceNewRelicCloudGcpDmLinkAccountRead,
		UpdateContext: resourceNewRelicCloudGcpDmLinkAccountUpdate,
		DeleteContext: resourceNewRelicCloudGcpDmLinkAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The New Relic account ID to link the GCP project to. Defaults to the provider account if not set.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name for this linked GCP account in New Relic.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The GCP project ID to link (e.g. 'my-gcp-project-123').",
			},
			"audience": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The Workload Identity Federation pool provider audience URI. " +
					"Format: //iam.googleapis.com/projects/{PROJECT_NUMBER}/locations/global/" +
					"workloadIdentityPools/{POOL_ID}/providers/{PROVIDER_ID}. " +
					"This is the 'name' attribute of the google_iam_workload_identity_pool_provider " +
					"resource prefixed with '//iam.googleapis.com/'.",
			},
			"service_account_email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The GCP service account email that New Relic will impersonate to collect metrics. " +
					"The service account must have the monitoring.viewer (and optionally " +
					"serviceusage.serviceUsageConsumer) role and must grant the WIF pool the " +
					"roles/iam.workloadIdentityUser binding.",
			},
		},
	}
}

// gcpDmOIDCEndpoint returns the New Relic OIDC token endpoint for the given provider region.
// This URL is set as credential_source.url in the WIF credential JSON and tells GCP STS
// where to fetch the subject token from.
func gcpDmOIDCEndpoint(region string) string {
	switch strings.ToLower(strings.TrimSpace(region)) {
	case "eu":
		return "https://oidc.eu.newrelic.com/r/gcp-cmp"
	case "staging":
		return "https://oidc-staging.newrelic.com/r/gcp-cmp"
	default: // US and JP use the US endpoint
		return "https://oidc.newrelic.com/r/gcp-cmp"
	}
}

// gcpDmBuildWIFCredential constructs the GCP Workload Identity Federation credential
// JSON string that cloudAuthenticateIntegration expects as its payload.
// All fixed fields (universe_domain, type, subject_token_type, token_url, format)
// are set to their required values; the caller supplies only the environment-specific inputs.
func gcpDmBuildWIFCredential(audience, serviceAccountEmail, region string) (string, error) {
	cred := map[string]interface{}{
		"universe_domain":    "googleapis.com",
		"type":               "external_account",
		"audience":           audience,
		"subject_token_type": "urn:ietf:params:oauth:token-type:jwt",
		"token_url":          "https://sts.googleapis.com/v1/token",
		"credential_source": map[string]interface{}{
			"url":     gcpDmOIDCEndpoint(region),
			"headers": map[string]interface{}{},
			"format": map[string]interface{}{
				"type":                    "json",
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

func resourceNewRelicCloudGcpDmLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	// Build the WIF credential JSON from individual fields.
	wifCredential, err := gcpDmBuildWIFCredential(
		d.Get("audience").(string),
		d.Get("service_account_email").(string),
		providerConfig.Region,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to build WIF credential: %w", err))
	}

	// Step 1: Authenticate via WIF to obtain an authReferenceId (30-min TTL).
	var authResp gcpDmAuthResp
	authVars := map[string]interface{}{
		"accountId":    accountID,
		"providerSlug": "GCP",
		"authType":     "WIF",
		"payload":      wifCredential,
	}
	if err := client.NerdGraph.QueryWithResponseAndContext(ctx, gcpDmAuthenticateMutation, authVars, &authResp); err != nil {
		return diag.FromErr(fmt.Errorf("cloudAuthenticateIntegration failed: %w", err))
	}
	authReferenceId := authResp.CloudAuthenticateIntegration.AuthReferenceId
	if authReferenceId == "" {
		return diag.FromErr(fmt.Errorf("cloudAuthenticateIntegration returned empty authReferenceId"))
	}

	// Step 2: Link GCP project to New Relic using the authReferenceId.
	var linkResp gcpDmLinkResp
	linkVars := map[string]interface{}{
		"accountId": accountID,
		"accounts": map[string]interface{}{
			"gcp": []gcpDmLinkAccountInput{{
				Name:            d.Get("name").(string),
				ProjectId:       d.Get("project_id").(string),
				AuthReferenceId: authReferenceId,
			}},
		},
	}
	if err := client.NerdGraph.QueryWithResponseAndContext(ctx, gcpDmLinkAccountMutation, linkVars, &linkResp); err != nil {
		return diag.FromErr(fmt.Errorf("cloudLinkAccount failed: %w", err))
	}

	if len(linkResp.CloudLinkAccount.Errors) > 0 {
		msgs := make([]string, 0, len(linkResp.CloudLinkAccount.Errors))
		for _, e := range linkResp.CloudLinkAccount.Errors {
			msgs = append(msgs, e.Type+": "+e.Message)
		}
		return diag.FromErr(fmt.Errorf("cloudLinkAccount errors: %s", strings.Join(msgs, "; ")))
	}

	if len(linkResp.CloudLinkAccount.LinkedAccounts) == 0 {
		return diag.FromErr(fmt.Errorf("cloudLinkAccount returned no linked accounts"))
	}

	d.SetId(strconv.Itoa(linkResp.CloudLinkAccount.LinkedAccounts[0].ID))
	_ = d.Set("account_id", accountID)

	return resourceNewRelicCloudGcpDmLinkAccountRead(ctx, d, meta)
}

func resourceNewRelicCloudGcpDmLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Use a minimal custom query that fetches only basic fields without integrations.
	// GetLinkedAccountWithContext uses client-go inline fragments that fail to resolve
	// GCP v2-specific integration types, causing "Abstract type must resolve" errors.
	var resp gcpDmLinkedAccountResp
	vars := map[string]interface{}{
		"accountId":       accountID,
		"linkedAccountId": linkedAccountID,
	}
	if err := client.NerdGraph.QueryWithResponseAndContext(ctx, gcpDmGetLinkedAccountQuery, vars, &resp); err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	la := resp.Actor.Account.Cloud.LinkedAccount
	if la == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("account_id", la.NrAccountId)
	_ = d.Set("name", la.Name)
	_ = d.Set("project_id", la.ExternalId)
	// audience and service_account_email are write-only (ForceNew); the API does not return them.
	// Terraform retains their values from state written during Create.

	return nil
}

func resourceNewRelicCloudGcpDmLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	renameInput := []cloud.CloudRenameAccountsInput{
		{
			LinkedAccountId: linkedAccountID,
			Name:            d.Get("name").(string),
		},
	}

	cloudRenamePayload, err := client.Cloud.CloudRenameAccountWithContext(ctx, accountID, renameInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudRenameAccount failed: %w", err))
	}

	var diags diag.Diagnostics
	if len(cloudRenamePayload.Errors) > 0 {
		for _, e := range cloudRenamePayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  e.Type + " " + e.Message,
			})
		}
		return diags
	}

	return resourceNewRelicCloudGcpDmLinkAccountRead(ctx, d, meta)
}

func resourceNewRelicCloudGcpDmLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	unlinkInput := []cloud.CloudUnlinkAccountsInput{
		{LinkedAccountId: linkedAccountID},
	}

	cloudUnlinkPayload, err := client.Cloud.CloudUnlinkAccountWithContext(ctx, accountID, unlinkInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudUnlinkAccount failed: %w", err))
	}

	var diags diag.Diagnostics
	if len(cloudUnlinkPayload.Errors) > 0 {
		for _, e := range cloudUnlinkPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  e.Type + " " + e.Message,
			})
		}
		return diags
	}

	d.SetId("")
	return nil
}
