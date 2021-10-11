package newrelic

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
)

func resourceNewRelicApplicationSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicApplicationSettingsCreate,
		ReadContext:   resourceNewRelicApplicationSettingsRead,
		UpdateContext: resourceNewRelicApplicationSettingsUpdate,
		DeleteContext: resourceNewRelicApplicationSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_apdex_threshold": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"end_user_apdex_threshold": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"enable_real_user_monitoring": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceNewRelicApplicationSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	userApp := expandApplication(d)

	listParams := apm.ListApplicationsParams{
		Name: userApp.Name,
	}

	result, err := client.APM.ListApplicationsWithContext(ctx, &listParams)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(result) != 1 {
		return diag.Errorf("more/less than one result from query for %s", userApp.Name)
	}

	app := *result[0]

	if app.Name != userApp.Name {
		return diag.Errorf("the result name %s does not match requested name %s", app.Name, userApp.Name)
	}

	d.SetId(strconv.Itoa(app.ID))

	log.Printf("[INFO] Importing New Relic application %v", userApp.Name)
	return resourceNewRelicApplicationSettingsUpdate(ctx, d, meta)
}

func resourceNewRelicApplicationSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	userApp := expandApplication(d)
	log.Printf("[INFO] Reading New Relic application %+v", userApp)

	app, err := client.APM.GetApplicationWithContext(ctx, userApp.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Read found New Relic application %+v\n\n\n", app)

	return diag.FromErr(flattenApplication(app, d))
}

func resourceNewRelicApplicationSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	userApp := expandApplication(d)

	updateParams := apm.UpdateApplicationParams{
		Name:     userApp.Name,
		Settings: userApp.Settings,
	}

	log.Printf("[INFO] Updating New Relic application %+v with params: %+v", userApp, updateParams)

	app, err := client.APM.UpdateApplicationWithContext(ctx, userApp.ID, updateParams)
	if err != nil {
		return diag.FromErr(err)
	}

	time.Sleep(2 * time.Second)

	return diag.FromErr(flattenApplication(app, d))
}

func resourceNewRelicApplicationSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// You can not delete application settings
	return nil
}
