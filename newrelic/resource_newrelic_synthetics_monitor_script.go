package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func resourceNewRelicSyntheticsMonitorScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsMonitorScriptCreate,
		ReadContext:   resourceNewRelicSyntheticsMonitorScriptRead,
		UpdateContext: resourceNewRelicSyntheticsMonitorScriptUpdate,
		DeleteContext: resourceNewRelicSyntheticsMonitorScriptDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSyntheticsMonitorScript,
		},
		Schema: map[string]*schema.Schema{
			"monitor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the monitor to attach the script to.",
			},
			"text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plaintext representing the monitor script.",
			},
		},
	}
}

func importSyntheticsMonitorScript(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	_ = d.Set("monitor_id", d.Id())
	return []*schema.ResourceData{d}, nil
}

func buildSyntheticsMonitorScriptStruct(d *schema.ResourceData) *synthetics.MonitorScript {
	script := synthetics.MonitorScript{
		Text: d.Get("text").(string),
	}

	return &script
}

func resourceNewRelicSyntheticsMonitorScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	id := d.Get("monitor_id").(string)
	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", id)

	_, err := client.Synthetics.UpdateMonitorScriptWithContext(ctx, id, *buildSyntheticsMonitorScriptStruct(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceNewRelicSyntheticsMonitorScriptRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics script %s", d.Id())

	script, err := client.Synthetics.GetMonitorScriptWithContext(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	_ = d.Set("text", script.Text)
	return nil
}

func resourceNewRelicSyntheticsMonitorScriptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", d.Id())

	_, err := client.Synthetics.UpdateMonitorScriptWithContext(ctx, d.Id(), *buildSyntheticsMonitorScriptStruct(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Id())
	return resourceNewRelicSyntheticsMonitorScriptRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic Synthetics monitor script %s", d.Id())

	script := synthetics.MonitorScript{
		Text: " ",
	}

	if _, err := client.Synthetics.UpdateMonitorScriptWithContext(ctx, d.Id(), script); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
