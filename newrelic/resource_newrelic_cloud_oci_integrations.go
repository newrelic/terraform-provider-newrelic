package newrelic

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func resourceNewrelicCloudOciIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewrelicCloudOciIntegrationsCreate,
		ReadContext:   resourceNewrelicCloudOciIntegrationsRead,
		UpdateContext: resourceNewrelicCloudOciIntegrationsUpdate,
		DeleteContext: resourceNewrelicCloudOciIntegrationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: generateOciIntegrationSchema(),
	}
}

func generateOciIntegrationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Description: "ID of the newrelic account",
			Computed:    true,
			Optional:    true,
		},
		"linked_account_id": {
			Type:        schema.TypeInt,
			Description: "Id of the linked OCI account in New Relic",
			Required:    true,
			ForceNew:    true,
		},
		"oci_metadata_and_tags": {
			Type:        schema.TypeList,
			Description: "OCI Metadata and Tags integration",
			Elem:        &schema.Resource{},
			Optional:    true,
			MaxItems:    1,
		},
	}
}

func resourceNewrelicCloudOciIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	cloudOciIntegrationinputs, _ := expandCloudOciIntegrationsinputs(d)
	ociIntegrationspayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloudOciIntegrationinputs)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	if len(ociIntegrationspayload.Errors) > 0 {
		for _, err := range ociIntegrationspayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	if len(ociIntegrationspayload.Integrations) > 0 {
		d.SetId(strconv.Itoa(d.Get("linked_account_id").(int)))
	}
	return nil
}

func expandCloudOciIntegrationsinputs(d *schema.ResourceData) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {
	ociCloudIntegrations := cloud.CloudOciIntegrationsInput{}
	ociDisableIntegrations := cloud.CloudOciDisableIntegrationsInput{}
	var linkedAccountID int
	if lid, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = lid.(int)
	}
	if v, ok := d.GetOk("oci_metadata_and_tags"); ok {
		ociCloudIntegrations.OciMetadataAndTags = expandCloudOciIntegrationsInputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("oci_metadata_and_tags"); len(n.([]interface{})) < len(o.([]interface{})) {
		ociDisableIntegrations.OciMetadataAndTags = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	configureInput := cloud.CloudIntegrationsInput{
		Oci: ociCloudIntegrations,
	}
	disableInput := cloud.CloudDisableIntegrationsInput{
		Oci: ociDisableIntegrations,
	}
	return configureInput, disableInput
}

func expandCloudOciIntegrationsInputs(b []interface{}, linkedAccountID int) []cloud.CloudOciMetadataAndTagsIntegrationInput {
	expanded := make([]cloud.CloudOciMetadataAndTagsIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudOciMetadataAndTagsIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		input.LinkedAccountId = linkedAccountID
		expanded[i] = input
	}
	return expanded
}

func resourceNewrelicCloudOciIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	flattenCloudOciLinkedAccount(d, linkedAccount)
	return nil
}

func flattenCloudOciLinkedAccount(d *schema.ResourceData, linkedAccount *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", linkedAccount.NrAccountId)
	_ = d.Set("linked_account_id", linkedAccount.ID)
}

func resourceNewrelicCloudOciIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	configureInput, disableInput := expandCloudOciIntegrationsinputs(d)
	cloudDisableIntegrationsPayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudDisableIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudDisableIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	cloudOciIntegrationsPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, configureInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if len(cloudOciIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudOciIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	return nil
}

func resourceNewrelicCloudOciIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	deleteInput := expandCloudOciDisableinputs(d)
	ociDisablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, deleteInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(ociDisablePayload.Errors) > 0 {
		for _, err := range ociDisablePayload.Errors {
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

func expandCloudOciDisableinputs(d *schema.ResourceData) cloud.CloudDisableIntegrationsInput {
	cloudOciDisableInput := cloud.CloudOciDisableIntegrationsInput{}
	var linkedAccountID int
	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if _, ok := d.GetOk("oci_metadata_and_tags"); ok {
		cloudOciDisableInput.OciMetadataAndTags = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	deleteInput := cloud.CloudDisableIntegrationsInput{
		Oci: cloudOciDisableInput,
	}
	return deleteInput
}
