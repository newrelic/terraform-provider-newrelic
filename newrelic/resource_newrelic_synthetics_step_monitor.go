package newrelic

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNewRelicSyntheticsStepMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsStepMonitorCreate,
		ReadContext:   resourceNewRelicSyntheticsStepMonitorRead,
		UpdateContext: resourceNewRelicSyntheticsStepMonitorUpdate,
		DeleteContext: resourceNewRelicSyntheticsStepMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "ID of the newrelic account",
				Computed:    true,
				Optional:    true,
			},
			"enable_screenshot_on_failure_and_script": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Capture a screenshot during job execution.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the monitor in New Relic.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of this monitor.",
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The monitor status (i.e. ENABLED, MUTED, DISABLED).",
				ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorStatuses(), false),
			},
			"locations_public": {
				Type:         schema.TypeList,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  "Publicly available location names in which the monitor will run.",
				Optional:     true,
				AtLeastOneOf: []string{"locations_public", "locations_private"},
			},
			"locations_private": {
				Type:         schema.TypeList,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  "List private location GUIDs for which the monitor will run.",
				Optional:     true,
				AtLeastOneOf: []string{"locations_public", "locations_private"},
			},
			"tag": {
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				Description: "The tags that will be associated with the monitor",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the tag key",
						},
						"values": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Values associated with the tag key",
						},
					},
				},
			},
			"steps": {
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				Description: "The steps that make up the script the monitor will run",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ordinal": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The position of the step within the script ranging from 1-100",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The metadata values related to the step. valid values are ASSERT_ELEMENT, ASSERT_MODAL, ASSERT_TEXT, ASSERT_TITLE, CLICK_ELEMENT, DISMISS_MODAL, DOUBLE_CLICK_ELEMENT, HOVER_ELEMENT, NAVIGATE, SECURE_TEXT_ENTRY, SELECT_ELEMENT, TEXT_ENTRY",
							//ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorTypes(), false),
						},
						"values": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The metadata values related to the step",
						},
					},
				},
			},
			"period": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.",
				ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorPeriods(), false),
			},
		},
	}
}

func resourceNewRelicSyntheticsStepMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNewRelicSyntheticsStepMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil

}
func resourceNewRelicSyntheticsStepMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil

}
func resourceNewRelicSyntheticsStepMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil

}
