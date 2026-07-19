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

func resourceNewRelicCloudGcpLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudGcpLinkAccountCreate,
		ReadContext:   resourceNewRelicCloudGcpLinkAccountRead,
		UpdateContext: resourceNewRelicCloudGcpLinkAccountUpdate,
		DeleteContext: resourceNewRelicCloudGcpLinkAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "accountID of newrelic account",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of the linked account",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "project id of the Gcp account",
				Required:    true,
				ForceNew:    true,
			},
			// The two fields below opt the resource into GCP Dimensional Metrics
			// (v2) keyless linking via Workload Identity Federation (WIF). When both
			// are set, the resource authenticates via WIF and links under the gcp_v2
			// provider; when both are absent, it links using the legacy service-account
			// key flow (no auth reference). They must be supplied together.
			"audience": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"service_account_email"},
				Description: "The Workload Identity Federation pool provider audience URI, used for GCP " +
					"Dimensional Metrics (keyless) linking. Format: //iam.googleapis.com/projects/" +
					"{PROJECT_NUMBER}/locations/global/workloadIdentityPools/{POOL_ID}/providers/{PROVIDER_ID}. " +
					"When set together with service_account_email, the account is linked via WIF instead " +
					"of a service-account key.",
			},
			"service_account_email": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"audience"},
				Description: "The GCP service account email New Relic impersonates to collect metrics when " +
					"linking via Workload Identity Federation (GCP Dimensional Metrics). The service account " +
					"must grant the WIF pool the roles/iam.workloadIdentityUser binding. Required together " +
					"with audience.",
			},
		},
	}
}

func resourceNewRelicCloudGcpLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	// GCP Dimensional Metrics (WIF/keyless) mode is selected when both audience and
	// service_account_email are provided; otherwise fall back to the legacy flow.
	if isGcpWIFMode(d) {
		return resourceNewRelicCloudGcpLinkAccountCreateWIF(ctx, d, providerConfig, accountID)
	}

	linkAccountInput := expandGcpCloudLinkAccountInput(d)

	var diags diag.Diagnostics

	//cloudLinkAccountWithContext func which links Gcp account with Newrelic
	//which returns payload and error
	cloudLinkAccountPayload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if cloudLinkAccountPayload == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] cloudLinkAccountPayload was nil"))
	}

	if len(cloudLinkAccountPayload.Errors) > 0 {
		for _, err := range cloudLinkAccountPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
	}

	if len(cloudLinkAccountPayload.LinkedAccounts) > 0 {
		d.SetId(strconv.Itoa(cloudLinkAccountPayload.LinkedAccounts[0].ID))
	}

	return diags
}

// isGcpWIFMode reports whether the resource is configured for GCP Dimensional
// Metrics (keyless) linking via Workload Identity Federation. Both audience and
// service_account_email are required together (enforced by RequiredWith), so
// checking audience is sufficient.
func isGcpWIFMode(d *schema.ResourceData) bool {
	return d.Get("audience").(string) != ""
}

// resourceNewRelicCloudGcpLinkAccountCreateWIF performs the two-step GCP Dimensional
// Metrics linking flow: authenticate via WIF to obtain a short-lived authReferenceId,
// then link the GCP project to New Relic using that reference.
func resourceNewRelicCloudGcpLinkAccountCreateWIF(ctx context.Context, d *schema.ResourceData, providerConfig *ProviderConfig, accountID int) diag.Diagnostics {
	client := providerConfig.NewClient

	wifCredential, err := gcpBuildWIFCredential(
		d.Get("audience").(string),
		d.Get("service_account_email").(string),
		providerConfig.Region,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to build WIF credential: %w", err))
	}

	// Step 1: Authenticate via WIF to obtain an authReferenceId (30-min TTL).
	authPayload, err := client.Cloud.CloudAuthenticateIntegrationWithContext(ctx, accountID, "GCP", "WIF", wifCredential)
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudAuthenticateIntegration failed: %w", err))
	}
	if authPayload.AuthReferenceId == "" {
		return diag.FromErr(fmt.Errorf("cloudAuthenticateIntegration returned empty authReferenceId"))
	}

	// Step 2: Link the GCP project to New Relic using the authReferenceId.
	linkPayload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, cloud.CloudLinkCloudAccountsInput{
		Gcp: []cloud.CloudGcpLinkAccountInput{{
			Name:            d.Get("name").(string),
			ProjectId:       d.Get("project_id").(string),
			AuthReferenceId: authPayload.AuthReferenceId,
		}},
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudLinkAccount failed: %w", err))
	}

	if len(linkPayload.Errors) > 0 {
		msgs := make([]string, 0, len(linkPayload.Errors))
		for _, e := range linkPayload.Errors {
			msgs = append(msgs, e.Type+" "+e.Message)
		}
		return diag.FromErr(fmt.Errorf("cloudLinkAccount errors: %s", strings.Join(msgs, "; ")))
	}

	if len(linkPayload.LinkedAccounts) == 0 {
		return diag.FromErr(fmt.Errorf("cloudLinkAccount returned no linked accounts"))
	}

	d.SetId(strconv.Itoa(linkPayload.LinkedAccounts[0].ID))
	_ = d.Set("account_id", accountID)

	return resourceNewRelicCloudGcpLinkAccountRead(ctx, d, providerConfig)
}

// expand function to extract inputs from the schema.
// Here it takes ResourceData as input and returns cloudLinkCloudAccountsInput.
func expandGcpCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {

	gcpAccount := cloud.CloudGcpLinkAccountInput{}

	if name, ok := d.GetOk("name"); ok {
		gcpAccount.Name = name.(string)
	}

	if projectID, ok := d.GetOk("project_id"); ok {
		gcpAccount.ProjectId = projectID.(string)
	}

	input := cloud.CloudLinkCloudAccountsInput{
		Gcp: []cloud.CloudGcpLinkAccountInput{gcpAccount},
	}

	return input

}

// gcpWIFOIDCEndpoint returns the New Relic OIDC token endpoint for the given provider region.
// This URL is set as credential_source.url in the WIF credential JSON and tells GCP STS
// where to fetch the subject token from.
func gcpWIFOIDCEndpoint(region string) string {
	switch strings.ToLower(strings.TrimSpace(region)) {
	case "eu":
		return "https://oidc.eu.newrelic.com/r/gcp-cmp"
	case "staging":
		return "https://oidc-staging.newrelic.com/r/gcp-cmp"
	default: // US and JP use the US endpoint
		return "https://oidc.newrelic.com/r/gcp-cmp"
	}
}

// gcpBuildWIFCredential constructs the GCP Workload Identity Federation credential
// JSON string that cloudAuthenticateIntegration expects as its payload.
// All fixed fields (universe_domain, type, subject_token_type, token_url, format)
// are set to their required values; the caller supplies only the environment-specific inputs.
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

// gcpWIFGetLinkedAccountQuery fetches only the basic fields of a linked account
// (id, name, nrAccountId, externalId) without requesting integrations.
// This avoids the "Abstract type 'Integration' must resolve to an Object type"
// error that occurs when GetLinkedAccountWithContext encounters GCP Dimensional
// Metrics-specific integration types that its inline fragments don't cover.
const gcpWIFGetLinkedAccountQuery = `query($accountId: Int!, $linkedAccountId: Int!) {
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

// gcpWIFLinkedAccountResp is the response type for gcpWIFGetLinkedAccountQuery.
type gcpWIFLinkedAccountResp struct {
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

func resourceNewRelicCloudGcpLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	// readWIFMinimal reads a linked account with a minimal query that omits the
	// integrations field. GetLinkedAccountWithContext requests integrations, which
	// throws "Abstract type 'Integration' must resolve to an Object type" on GCP
	// Dimensional Metrics (gcp_v2) accounts. audience and service_account_email are
	// write-only (ForceNew) and never returned by the API, so they are retained from
	// state written during Create.
	readWIFMinimal := func() diag.Diagnostics {
		var resp gcpWIFLinkedAccountResp
		vars := map[string]interface{}{
			"accountId":       accountID,
			"linkedAccountId": linkedAccountID,
		}
		if err := client.NerdGraph.QueryWithResponseAndContext(ctx, gcpWIFGetLinkedAccountQuery, vars, &resp); err != nil {
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

		return nil
	}

	// WIF (GCP Dimensional Metrics) mode always uses the minimal query.
	if isGcpWIFMode(d) {
		return readWIFMinimal()
	}

	// Legacy (service-account-key) mode keeps the original read path unchanged.
	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		// A WIF-linked account imported without audience in state (mode is unknown at
		// import time) lands here; retry with the minimal query rather than failing.
		if strings.Contains(err.Error(), "must resolve to an Object type") {
			return readWIFMinimal()
		}
		return diag.FromErr(err)
	}

	readGcpLinkedAccount(d, linkedAccount)

	return nil

}

// readGcpLinkedAccount function to store name and ExternalId.
// Using set func to store the output values.
func readGcpLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("name", result.Name)
	_ = d.Set("project_id", result.ExternalId)
}

func resourceNewRelicCloudGcpLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	id, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	input := []cloud.CloudRenameAccountsInput{
		{
			Name:            d.Get("name").(string),
			LinkedAccountId: id,
		},
	}

	//CloudRenameAccount to rename the name of linkedAccount
	cloudRenameAccountPayload, err := client.Cloud.CloudRenameAccountWithContext(ctx, accountID, input)

	if err != nil {
		diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudRenameAccountPayload.Errors) > 0 {
		for _, err := range cloudRenameAccountPayload.Errors {

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})

		}

		return diags

	}

	return nil
}

func resourceNewRelicCloudGcpLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		diag.FromErr(convErr)
	}

	unlinkAccountInput := []cloud.CloudUnlinkAccountsInput{
		{
			LinkedAccountId: linkedAccountID,
		},
	}

	//CloudUnlinkAccountWithContext func to unlink the GCP account with Newrelic
	cloudUnlinkAccountPayload, err := client.Cloud.CloudUnlinkAccountWithContext(ctx, accountID, unlinkAccountInput)

	if err != nil {
		diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudUnlinkAccountPayload.Errors) > 0 {
		for _, err := range cloudUnlinkAccountPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}

		return diags

	}
	//Setting up the linked account id to null after destroying the resource.
	d.SetId("")

	return nil

}
