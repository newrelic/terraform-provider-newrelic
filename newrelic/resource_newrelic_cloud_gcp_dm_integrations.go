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

// Auxiliary code (the service table, schema, expand, and flatten) lives in
// structures_newrelic_cloud_gcp_dm_integrations.go; this file holds the resource
// definition and its CRUD operations.

func resourceNewrelicCloudGcpDmIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewrelicCloudGcpDmIntegrationsCreate,
		ReadContext:   resourceNewrelicCloudGcpDmIntegrationsRead,
		UpdateContext: resourceNewrelicCloudGcpDmIntegrationsUpdate,
		DeleteContext: resourceNewrelicCloudGcpDmIntegrationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: generateGcpDmIntegrationSchema(),
	}
}

func resourceNewrelicCloudGcpDmIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkedAccountID := d.Get("linked_account_id").(int)

	configureInput, _ := expandCloudGcpDmIntegrationsInput(d, linkedAccountID)

	configPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloud.CloudIntegrationsInput{Gcp: configureInput})
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudConfigureIntegration failed: %w", err))
	}
	if len(configPayload.Errors) > 0 {
		return diag.FromErr(gcpDmMutationErrors("cloudConfigureIntegration", configPayload.Errors))
	}

	d.SetId(strconv.Itoa(linkedAccountID))
	_ = d.Set("account_id", accountID)

	return resourceNewrelicCloudGcpDmIntegrationsRead(ctx, d, meta)
}

func resourceNewrelicCloudGcpDmIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

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

	_ = d.Set("account_id", linkedAccount.NrAccountId)
	_ = d.Set("linked_account_id", linkedAccount.ID)
	flattenGcpDmIntegrations(d, linkedAccount)

	return nil
}

func resourceNewrelicCloudGcpDmIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	configureInput, disableInput := expandCloudGcpDmIntegrationsInput(d, linkedAccountID)

	// Disable the services whose blocks were removed from the config, then
	// enable/update the services whose blocks are present.
	disablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, cloud.CloudDisableIntegrationsInput{Gcp: disableInput})
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudDisableIntegration failed: %w", err))
	}
	if filterErr := gcpDmFilterDisableErrors(disablePayload.Errors); filterErr != nil {
		return diag.FromErr(filterErr)
	}

	configPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloud.CloudIntegrationsInput{Gcp: configureInput})
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudConfigureIntegration failed: %w", err))
	}
	if len(configPayload.Errors) > 0 {
		return diag.FromErr(gcpDmMutationErrors("cloudConfigureIntegration", configPayload.Errors))
	}

	return resourceNewrelicCloudGcpDmIntegrationsRead(ctx, d, meta)
}

func resourceNewrelicCloudGcpDmIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, disableInput := expandCloudGcpDmIntegrationsInput(d, linkedAccountID)

	disablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, cloud.CloudDisableIntegrationsInput{Gcp: disableInput})
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudDisableIntegration failed: %w", err))
	}
	if filterErr := gcpDmFilterDisableErrors(disablePayload.Errors); filterErr != nil {
		return diag.FromErr(filterErr)
	}

	d.SetId("")
	return nil
}

// gcpDmMutationErrors aggregates cloud mutation errors into a single error.
func gcpDmMutationErrors(mutation string, errors []cloud.CloudIntegrationMutationError) error {
	errMessages := make([]string, 0, len(errors))
	for _, e := range errors {
		errMessages = append(errMessages, e.Type+": "+e.Message)
	}
	return fmt.Errorf("%s errors:\n %s", mutation, strings.Join(errMessages, "\n "))
}
