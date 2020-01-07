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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"text": {
				Type:     schema.TypeString,
				Required: true,
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
	client := meta.(*ProviderConfig).Synthetics
	script := buildSyntheticsMonitorScriptStruct(d)

	id := d.Get("monitor_id").(string)
	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", id)

	err := client.UpdateMonitorScript(id, *script)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceNewRelicSyntheticsMonitorScriptRead(d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	log.Printf("[INFO] Reading New Relic Synthetics script %s", d.Id())

	script, err := client.GetMonitorScript(d.Id())
	if err != nil {
		if _, ok := err.(*errors.ErrorNotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("text", script.Text)
	return nil
}

func resourceNewRelicSyntheticsMonitorScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics
	script := buildSyntheticsMonitorScriptStruct(d)

	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", d.Id())

	err := client.UpdateMonitorScript(d.Id(), *script)
	if err != nil {
		return err
	}

	d.SetId(d.Id())
	return resourceNewRelicSyntheticsMonitorScriptRead(d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	log.Printf("[INFO] Deleting New Relic Synthetics monitor script %s", d.Id())

	script := synthetics.MonitorScript{
		Text: " ",
	}

	if err := client.UpdateMonitorScript(d.Id(), script); err != nil {
		return err
	}

	return nil
}
