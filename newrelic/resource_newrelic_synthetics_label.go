package newrelic

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func resourceNewRelicSyntheticsLabel() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "`newrelic_synthetics_label` has been deprecated.  Use `newrelic_entity_tags` instead.",
		Create:             resourceNewRelicSyntheticsLabelCreate,
		Read:               resourceNewRelicSyntheticsLabelRead,
		Delete:             resourceNewRelicSyntheticsLabelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"monitor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the monitor that will be assigned the label.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A string representing the label key/category.",
				// Case fold this attribute when diffing
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A string representing the label value.",
				// Case fold this attribute when diffing
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the Synthetics label.",
			},
		},
	}
}

func resourceNewRelicSyntheticsLabelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	monitorID := d.Get("monitor_id").(string)
	label := expandSyntheticsLabel(d)

	log.Printf("[INFO] Creating New Relic Synthetics label %s:%s", label.Type, label.Value)

	// Note: AddMonitorLabel is deprecated
	err := client.Synthetics.AddMonitorLabel(monitorID, label.Type, label.Value)
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{monitorID, label.Type, label.Value}, ":"))

	return nil
}

func resourceNewRelicSyntheticsLabelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics label %s", d.Id())

	ids := strings.Split(d.Id(), ":")
	monitorID := ids[0]
	labelType := ids[1]
	value := ids[2]

	_, err := client.Synthetics.GetMonitor(monitorID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	// Note: GetMonitorLabels is deprecated
	labels, err := client.Synthetics.GetMonitorLabels(monitorID)
	if err != nil {
		return err
	}

	var label *synthetics.MonitorLabel
	for _, l := range labels {
		if !strings.EqualFold(l.Type, labelType) || !strings.EqualFold(l.Value, value) {
			continue
		}

		label = l
	}

	if label == nil {
		d.SetId("")
		return nil
	}

	return flattenSyntheticsLabel(label, d)
}

func resourceNewRelicSyntheticsLabelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids := strings.Split(d.Id(), ":")
	monitorID := ids[0]
	labelType := ids[1]
	value := ids[2]

	log.Printf("[INFO] Deleting New Relic alert condition %s", d.Id())

	// Note: DeleteMonitorLabel is deprecated
	err := client.Synthetics.DeleteMonitorLabel(monitorID, labelType, value)
	if err != nil {
		return err
	}

	return nil
}

func expandSyntheticsLabel(d *schema.ResourceData) *synthetics.MonitorLabel {
	label := synthetics.MonitorLabel{
		Type:  d.Get("type").(string),
		Value: d.Get("value").(string),
		Href:  d.Get("href").(string),
	}

	return &label
}

func flattenSyntheticsLabel(label *synthetics.MonitorLabel, d *schema.ResourceData) error {
	ids := strings.Split(d.Id(), ":")
	monitorID := ids[0]

	d.Set("monitor_id", monitorID)
	d.Set("type", label.Type)
	d.Set("value", label.Value)
	d.Set("href", label.Href)

	return nil
}
