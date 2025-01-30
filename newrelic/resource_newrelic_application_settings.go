package newrelic

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apm"
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

	var app *apm.Application

	for _, a := range result {
		if a.Name == userApp.Name {
			app = a
			break
		}
	}

	if app == nil {
		return diag.Errorf("the name '%s' does not match any New Relic applications", userApp.Name)
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
