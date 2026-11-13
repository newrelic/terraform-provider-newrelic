package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

func resourceNewRelicMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicMaintenanceWindowCreate,
		ReadContext:   resourceNewRelicMaintenanceWindowRead,
		UpdateContext: resourceNewRelicMaintenanceWindowUpdate,
		DeleteContext: resourceNewRelicMaintenanceWindowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the maintenance window.",
				ValidateFunc: validation.All(
					validation.StringIsNotWhiteSpace,
					validation.StringLenBetween(1, 100),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The description of the maintenance window.",
				ValidateFunc: validation.StringLenBetween(0, 350),
			},
			"scope": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "The scope the maintenance window belongs to.",
				Elem:        maintenanceWindowScopeSchema(),
			},
			"start_time": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The start time of the maintenance window, as an ISO 8601 date-time without a time zone (e.g. 2026-11-05T21:00:00).",
				ValidateFunc: validateMaintenanceWindowStartTime,
			},
			"duration": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The duration of the maintenance window, in ISO 8601 duration format (e.g. PT1H for one hour). Capped at 24 hours when `rrule` is set.",
				ValidateFunc: validateMaintenanceWindowDuration,
			},
			"rrule": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The recurrence rule of the maintenance window, in RRULE format, for recurring maintenance windows.",
				ValidateFunc: validateMaintenanceWindowRrule,
			},
			"timezone": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The timezone in which the maintenance window is scheduled (e.g. America/Los_Angeles).",
				ValidateFunc: validateMaintenanceWindowTimeZone,
			},
			"affected_entity_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The type of the entities affected by the maintenance window.",
				ValidateFunc: validation.StringInSlice([]string{"SERVICE_LEVEL"}, false),
			},
			"affected_entities": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    100,
				Description: "The list of entity IDs affected by the maintenance window.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		CustomizeDiff: validateMaintenanceWindowAttributes,
	}
}

func maintenanceWindowScopeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the object within the given scope (e.g. an account ID).",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The kind of scope `id` refers to. Valid values are `ACCOUNT` and `ORGANIZATION`.",
				ValidateFunc: validation.StringInSlice([]string{"ACCOUNT", "ORGANIZATION"}, false),
			},
		},
	}
}

func resourceNewRelicMaintenanceWindowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	createInput := expandMaintenanceWindowCreateInput(d)

	log.Printf("[INFO] Creating New Relic Maintenance Window %s", createInput.Name)

	created, err := client.MaintenanceWindow.MaintenanceWindowCreateWithContext(ctx, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.ID)

	return diag.FromErr(flattenMaintenanceWindow(created, d))
}

func resourceNewRelicMaintenanceWindowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	resp, err := client.MaintenanceWindow.GetListByIdsWithContext(ctx, []string{d.Id()})
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if resp == nil || len(resp.MaintenanceWindows) == 0 {
		d.SetId("")
		return nil
	}

	mw := resp.MaintenanceWindows[0]
	return diag.FromErr(flattenMaintenanceWindow(&mw, d))
}

func resourceNewRelicMaintenanceWindowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	updateInput := expandMaintenanceWindowUpdateInput(d)

	log.Printf("[INFO] Updating New Relic Maintenance Window %s", d.Id())

	updated, err := client.MaintenanceWindow.MaintenanceWindowUpdateWithContext(ctx, d.Id(), updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenMaintenanceWindow(updated, d))
}

func resourceNewRelicMaintenanceWindowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic Maintenance Window %s", d.Id())

	if _, err := client.MaintenanceWindow.MaintenanceWindowDeleteWithContext(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
