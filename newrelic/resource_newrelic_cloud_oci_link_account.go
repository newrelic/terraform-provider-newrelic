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

func resourceNewRelicCloudOciAccountLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudOciAccountLinkCreate,
		ReadContext:   resourceNewRelicCloudOciAccountLinkRead,
		UpdateContext: resourceNewRelicCloudOciAccountLinkUpdate,
		DeleteContext: resourceNewRelicCloudOciAccountLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID where you want to link the OCI account.",
				ForceNew:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "The OCI tenant identifier.",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The linked account name.",
				Required:    true,
			},
			"compartment_ocid": {
				Type:        schema.TypeString,
				Description: "The New Relic compartment OCID in OCI.",
				Required:    true,
			},
			"oci_client_id": {
				Type:        schema.TypeString,
				Description: "The client ID for OCI WIF.",
				Required:    true,
			},
			"oci_client_secret": {
				Type:        schema.TypeString,
				Description: "The client secret for OCI WIF.",
				Required:    true,
				Sensitive:   true,
			},
			"oci_domain_url": {
				Type:        schema.TypeString,
				Description: "The OCI domain URL for WIF.",
				Required:    true,
			},
			"oci_home_region": {
				Type:        schema.TypeString,
				Description: "The home region of the tenancy.",
				Required:    true,
			},
			"oci_svc_user_name": {
				Type:        schema.TypeString,
				Description: "The service user name for OCI WIF.",
				Required:    true,
			},
			"oci_region": {
				Type:        schema.TypeString,
				Description: "The OCI region for the account. This field is only used for updates, not during initial creation.",
				Optional:    true,
			},
			"metric_stack_ocid": {
				Type:        schema.TypeString,
				Description: "The metric stack identifier for the OCI account. This field is only used for updates, not during initial creation.",
				Optional:    true,
			},
			"ingest_vault_ocid": {
				Type:        schema.TypeString,
				Description: "The OCI ingest secret OCID.",
				Optional:    true,
			},
			"instrumentation_type": {
				Type:        schema.TypeString,
				Description: "Specifies the type of integration, such as metrics, logs, or a combination of logs and metrics.",
				Optional:    true,
			},
			"logging_stack_ocid": {
				Type:        schema.TypeString,
				Description: "The Logging stack identifier for the OCI account.",
				Optional:    true,
			},
			"user_vault_ocid": {
				Type:        schema.TypeString,
				Description: "The user secret OCID.",
				Optional:    true,
			},
		},
	}
}

func resourceNewRelicCloudOciAccountLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkAccountInput := expandOciCloudLinkAccountInput(d)

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

func expandOciCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	ociAccount := cloud.CloudOciLinkAccountInput{}

	if tenantID, ok := d.GetOk("tenant_id"); ok {
		ociAccount.TenantId = tenantID.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		ociAccount.Name = name.(string)
	}

	if compartmentOcid, ok := d.GetOk("compartment_ocid"); ok {
		ociAccount.CompartmentOcid = compartmentOcid.(string)
	}

	if ociClientID, ok := d.GetOk("oci_client_id"); ok {
		ociAccount.OciClientId = ociClientID.(string)
	}

	if ociClientSecret, ok := d.GetOk("oci_client_secret"); ok {
		ociAccount.OciClientSecret = cloud.SecureValue(ociClientSecret.(string))
	}

	if ociDomainURL, ok := d.GetOk("oci_domain_url"); ok {
		ociAccount.OciDomainURL = ociDomainURL.(string)
	}

	if ociHomeRegion, ok := d.GetOk("oci_home_region"); ok {
		ociAccount.OciHomeRegion = ociHomeRegion.(string)
	}

	if ingestVaultOcid, ok := d.GetOk("ingest_vault_ocid"); ok {
		ociAccount.IngestVaultOcid = ingestVaultOcid.(string)
	}

	if instrumentationType, ok := d.GetOk("instrumentation_type"); ok {
		ociAccount.InstrumentationType = instrumentationType.(string)
	}

	if userVaultOcid, ok := d.GetOk("user_vault_ocid"); ok {
		ociAccount.UserVaultOcid = userVaultOcid.(string)
	}

	input := cloud.CloudLinkCloudAccountsInput{
		Oci: []cloud.CloudOciLinkAccountInput{ociAccount},
	}
	return input
}

func resourceNewRelicCloudOciAccountLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	readOciLinkedAccount(d, linkedAccount)

	return nil
}

func readOciLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("tenant_id", result.ExternalId)
	_ = d.Set("name", result.Name)
}

func expandOciCloudUpdateAccountInput(d *schema.ResourceData) cloud.CloudUpdateCloudAccountsInput {
	linkedAccountID, _ := strconv.Atoi(d.Id())

	ociAccount := cloud.CloudOciUpdateAccountInput{
		LinkedAccountId: linkedAccountID,
	}

	if name, ok := d.GetOk("name"); ok {
		ociAccount.Name = name.(string)
	}

	if compartmentOcid, ok := d.GetOk("compartment_ocid"); ok {
		ociAccount.CompartmentOcid = compartmentOcid.(string)
	}

	if ociClientID, ok := d.GetOk("oci_client_id"); ok {
		ociAccount.OciClientId = ociClientID.(string)
	}

	if ociClientSecret, ok := d.GetOk("oci_client_secret"); ok {
		ociAccount.OciClientSecret = cloud.SecureValue(ociClientSecret.(string))
	}

	if ociDomainURL, ok := d.GetOk("oci_domain_url"); ok {
		ociAccount.OciDomainURL = ociDomainURL.(string)
	}

	if ociHomeRegion, ok := d.GetOk("oci_home_region"); ok {
		ociAccount.OciHomeRegion = ociHomeRegion.(string)
	}

	if tenantID, ok := d.GetOk("tenant_id"); ok {
		ociAccount.TenantId = tenantID.(string)
	}

	if ociRegion, ok := d.GetOk("oci_region"); ok {
		ociAccount.OciRegion = ociRegion.(string)
	}

	if metricStackOcid, ok := d.GetOk("metric_stack_ocid"); ok {
		ociAccount.MetricStackOcid = metricStackOcid.(string)
	}

	if ingestVaultOcid, ok := d.GetOk("ingest_vault_ocid"); ok {
		ociAccount.IngestVaultOcid = ingestVaultOcid.(string)
	}

	if instrumentationType, ok := d.GetOk("instrumentation_type"); ok {
		ociAccount.InstrumentationType = instrumentationType.(string)
	}

	if loggingStackOcid, ok := d.GetOk("logging_stack_ocid"); ok {
		ociAccount.LoggingStackOcid = loggingStackOcid.(string)
	}

	if userVaultOcid, ok := d.GetOk("user_vault_ocid"); ok {
		ociAccount.UserVaultOcid = userVaultOcid.(string)
	}

	input := cloud.CloudUpdateCloudAccountsInput{
		Oci: []cloud.CloudOciUpdateAccountInput{ociAccount},
	}
	return input
}

func resourceNewRelicCloudOciAccountLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	updateAccountInput := expandOciCloudUpdateAccountInput(d)

	cloudUpdateAccountPayload, err := client.Cloud.CloudUpdateAccountWithContext(ctx, accountID, updateAccountInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(cloudUpdateAccountPayload.LinkedAccounts) == 0 {
		linkedAccountID, _ := strconv.Atoi(d.Id())
		return diag.FromErr(fmt.Errorf("no linked account with 'linked_account_id': %d found", linkedAccountID))
	}

	return nil
}

func resourceNewRelicCloudOciAccountLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	unlinkAccountInput := []cloud.CloudUnlinkAccountsInput{
		{
			LinkedAccountId: linkedAccountID,
		},
	}
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

	d.SetId("")
	return nil
}
