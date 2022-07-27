package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Schema: mergeSchemas(
			syntheticsMonitorCommonSchema(),
			syntheticsStepMonitorSchema(),
		),
	}
}

func syntheticsStepMonitorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_screenshot_on_failure": {
			Type:        schema.TypeBool,
			Description: "Capture a screenshot during job execution.",
			Optional:    true,
		},
		"location_private": {
			Type:        schema.TypeSet,
			Description: "",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"guid": {
						Type:        schema.TypeString,
						Description: "The unique identifier for the Synthetics private location in New Relic.",
						Required:    true,
					},
					"vse_password": {
						Type:        schema.TypeString,
						Description: "The location's Verified Script Execution password (Only necessary if Verified Script Execution is enabled for the location).",
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
		},
		"location_public": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Optional:    true,
			Description: "The public location(s) that the monitor will run jobs from.",
		},
		"steps": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "The steps that make up the script the monitor will run",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ordinal": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "The position of the step within the script ranging from 0-100",
						Default:     0,
						//SchemaValidateDiagFunc: // TODO: add validation to ensure value is between 0 and 100 (inclusive)
					},
					"type": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The type of step to be added to the script.",
						//SchemaValidateDiagFunc: // TODO: add valid step types via enum values in client
					},
					"values": {
						Type:        schema.TypeList,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Optional:    true,
						Description: "",
					},
				},
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
