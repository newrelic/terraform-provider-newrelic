package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicOneDashboardJson() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicOneDashboardJsonCreate,
		ReadContext:   resourceNewRelicOneDashboardJsonRead,
		UpdateContext: resourceNewRelicOneDashboardJsonUpdate,
		DeleteContext: resourceNewRelicOneDashboardJsonDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dashboard's json.",
			},
			// Optional
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID where you want to create the dashboard.",
			},
		},
	}
}

func resourceNewRelicOneDashboardJsonCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Create")
	}

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	defaultInfo := map[string]interface{}{
		"account_id": accountID,
	}
	dashboard, err := expandDashboardJsonInput(d, defaultInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic One JSON dashboard: %s", dashboard.Name)

	created, err := client.Dashboards.DashboardCreateWithContext(ctx, accountID, *dashboard)
	if err != nil {
		return diag.FromErr(err)
	}
	guid := created.EntityResult.GUID
	if guid == "" {
		var errMessages string
		for _, e := range created.Errors {
			errMessages += "[" + string(e.Type) + ": " + e.Description + "]"
		}

		return diag.Errorf("err: newrelic_one_dashboard_json Create failed: %s", errMessages)
	}

	log.Printf("[INFO] New JSON Dashboard GUID: %s", guid)

	d.SetId(string(guid))

	res := resourceNewRelicOneDashboardJsonRead(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return res
}

// resourceNewRelicOneDashboardRead NerdGraph => Terraform reader
func resourceNewRelicOneDashboardJsonRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic One JSON dashboard %s", d.Id())

	dashboard, err := client.Dashboards.GetDashboardEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenDashboardJsonEntity(dashboard, d))
}

func resourceNewRelicOneDashboardJsonUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	defaultInfo := map[string]interface{}{
		"account_id": accountID,
	}

	dashboard, err := expandDashboardJsonInput(d, defaultInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating New Relic One JSON dashboard '%s' (%s)", dashboard.Name, d.Id())

	updated, err := client.Dashboards.DashboardUpdateWithContext(ctx, *dashboard, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	guid := updated.EntityResult.GUID
	if guid == "" {
		var errMessages string
		for _, e := range updated.Errors {
			errMessages += "[" + string(e.Type) + ": " + e.Description + "]"
		}

		return diag.Errorf("err: newrelic_one_dashboard_json Update failed: %s", errMessages)
	}

	return resourceNewRelicOneDashboardJsonRead(ctx, d, meta)
}

func resourceNewRelicOneDashboardJsonDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One JSON dashboard %v", d.Id())

	if _, err := client.Dashboards.DashboardDeleteWithContext(ctx, common.EntityGUID(d.Id())); err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
