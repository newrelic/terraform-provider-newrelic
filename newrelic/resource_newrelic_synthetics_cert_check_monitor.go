package newrelic

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNewRelicSyntheticsCertCheckMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsCertCheckMonitorCreate,
		ReadContext:   resourceNewRelicSyntheticsCertCheckMonitorRead,
		UpdateContext: resourceNewRelicSyntheticsCertCheckMonitorUpdate,
		DeleteContext: resourceNewRelicSyntheticsCertCheckMonitorDelete,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "ID of the newrelic account",
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of the cert check monitor",
				Required:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "",
				Required:    true,
			},
			"certificate_expiration": {
				Type:        schema.TypeInt,
				Description: "",
				Required:    true,
			},
			"location_public": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				MinItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"location_public", "location_private"},
				Description:  "The locations in which this monitor should be run.",
			},
			"location_private": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				MinItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"location_public", "location_private"},
				Description:  "The locations in which this monitor should be run.",
			},
			"status": {
				Type:         schema.TypeString,
				Description:  "The monitor status (i.e. ENABLED, MUTED, DISABLED).",
				Required:     true,
				ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorStatuses(), false),
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
			"period": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.",
				ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorPeriods(), false),
			},
		},
	}
}

func resourceNewRelicSyntheticsCertCheckMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicSyntheticsCertCheckMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicSyntheticsCertCheckMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicSyntheticsCertCheckMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
