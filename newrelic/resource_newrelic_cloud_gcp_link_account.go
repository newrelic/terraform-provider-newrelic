package newrelic

import (
	"context"
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
		CustomizeDiff: resourceNewRelicCloudGcpLinkAccountCustomizeDiff,
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
			// use_workload_identity_federation explicitly selects GCP Dimensional Metrics (v2)
			// keyless linking via Workload Identity Federation (WIF). When true, the
			// resource authenticates via WIF and links under the gcp_v2 provider, and
			// audience + service_account_email are required. When false (the default),
			// the account is linked using the legacy service-account-key flow.
			"use_workload_identity_federation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				Description: "Set to true to link this GCP account for New Relic GCP Dimensional Metrics (v2) " +
					"using keyless Workload Identity Federation (WIF). When true, audience and " +
					"service_account_email are required. When false (the default), the account is linked " +
					"using the legacy service-account-key flow.",
			},
			"audience": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"service_account_email"},
				Description: "The Workload Identity Federation pool provider audience URI, used for GCP " +
					"Dimensional Metrics (keyless) linking. Format: //iam.googleapis.com/projects/" +
					"{PROJECT_NUMBER}/locations/global/workloadIdentityPools/{POOL_ID}/providers/{PROVIDER_ID}. " +
					"Required when use_workload_identity_federation = true.",
			},
			"service_account_email": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"audience"},
				Description: "The GCP service account email New Relic impersonates to collect metrics when " +
					"linking via Workload Identity Federation (GCP Dimensional Metrics). The service account " +
					"must grant the WIF pool the roles/iam.workloadIdentityUser binding. Required when " +
					"use_workload_identity_federation = true.",
			},
		},
	}
}

func resourceNewRelicCloudGcpLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	// GCP Dimensional Metrics (WIF/keyless) mode is selected by use_workload_identity_federation;
	// otherwise fall back to the legacy service-account-key flow.
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
// Metrics (v2) keyless linking via Workload Identity Federation, as selected by
// the explicit use_workload_identity_federation attribute.
func isGcpWIFMode(d *schema.ResourceData) bool {
	return d.Get("use_workload_identity_federation").(bool)
}

// resourceNewRelicCloudGcpLinkAccountCustomizeDiff enforces the relationship between
// the explicit use_workload_identity_federation flag and the WIF credential fields: when the
// flag is set, audience and service_account_email are required; when it is not set,
// those fields must be absent.
func resourceNewRelicCloudGcpLinkAccountCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	isDM := d.Get("use_workload_identity_federation").(bool)

	// audience and service_account_email are commonly derived from google_* resources
	// created in the same apply, so their resolved values are unknown at plan time
	// (d.Get would return ""). Validate presence from the raw config instead, which
	// reports an attribute as set even when its value is not yet known.
	audienceSet, saSet := false, false
	if rawConfig := d.GetRawConfig(); !rawConfig.IsNull() {
		audienceSet = !rawConfig.GetAttr("audience").IsNull()
		saSet = !rawConfig.GetAttr("service_account_email").IsNull()
	}

	if isDM && (!audienceSet || !saSet) {
		return fmt.Errorf("`audience` and `service_account_email` are required when `use_workload_identity_federation = true`")
	}
	if !isDM && (audienceSet || saSet) {
		return fmt.Errorf("`audience` and `service_account_email` can only be set when `use_workload_identity_federation = true`")
	}
	return nil
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
		errMessages := make([]string, 0, len(linkPayload.Errors))
		for _, e := range linkPayload.Errors {
			errMessages = append(errMessages, e.Type+" "+e.Message)
		}
		return diag.FromErr(fmt.Errorf("cloudLinkAccount errors:\n %s", strings.Join(errMessages, "\n ")))
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

func resourceNewRelicCloudGcpLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	// Both WIF (gcp_v2) and legacy accounts are read through the Go Client. For
	// WIF-linked accounts, audience and service_account_email are write-only
	// (ForceNew) and never returned by the API, so they are retained from state.
	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if linkedAccount == nil {
		d.SetId("")
		return nil
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
