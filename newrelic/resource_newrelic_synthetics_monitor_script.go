package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func resourceNewRelicSyntheticsMonitorScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicSyntheticsMonitorScriptCreate,
		Read:   resourceNewRelicSyntheticsMonitorScriptRead,
		Update: resourceNewRelicSyntheticsMonitorScriptUpdate,
		Delete: resourceNewRelicSyntheticsMonitorScriptDelete,
		Importer: &schema.ResourceImporter{
			State: importSyntheticsMonitorScript,
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

func importSyntheticsMonitorScript(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	d.Set("monitor_id", d.Id())
	return []*schema.ResourceData{d}, nil
}

func buildSyntheticsMonitorScriptStruct(d *schema.ResourceData) *synthetics.MonitorScript {
	script := synthetics.MonitorScript{
		Text: d.Get("text").(string),
	}

	return &script
}

func resourceNewRelicSyntheticsMonitorScriptCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	id := d.Get("monitor_id").(string)
	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", id)

	_, err := client.Synthetics.UpdateMonitorScript(id, *buildSyntheticsMonitorScriptStruct(d))
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceNewRelicSyntheticsMonitorScriptRead(d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics script %s", d.Id())

	script, err := client.Synthetics.GetMonitorScript(d.Id())
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("text", script.Text)
	return nil
}

func resourceNewRelicSyntheticsMonitorScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", d.Id())

	_, err := client.Synthetics.UpdateMonitorScript(d.Id(), *buildSyntheticsMonitorScriptStruct(d))
	if err != nil {
		return err
	}

	d.SetId(d.Id())
	return resourceNewRelicSyntheticsMonitorScriptRead(d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic Synthetics monitor script %s", d.Id())

	script := synthetics.MonitorScript{
		Text: " ",
	}

	if _, err := client.Synthetics.UpdateMonitorScript(d.Id(), script); err != nil {
		return err
	}

	return nil
}
