package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicSyntheticsAlertCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicSyntheticsAlertConditionCreate,
		Read:   resourceNewRelicSyntheticsAlertConditionRead,
		Update: resourceNewRelicSyntheticsAlertConditionUpdate,
		Delete: resourceNewRelicSyntheticsAlertConditionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"monitor_id": {
				Type:     schema.TypeString,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"runbook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func expandSyntheticsCondition(d *schema.ResourceData) *alerts.SyntheticsCondition {
	condition := alerts.SyntheticsCondition{
		Name:      d.Get("name").(string),
		Enabled:   d.Get("enabled").(bool),
		MonitorID: d.Get("monitor_id").(string),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	return &condition
}

func flattenSyntheticsCondition(condition *alerts.SyntheticsCondition, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	d.Set("policy_id", policyID)
	d.Set("monitor_id", condition.MonitorID)
	d.Set("name", condition.Name)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("enabled", condition.Enabled)

	return nil
}

func resourceNewRelicSyntheticsAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	policyID := d.Get("policy_id").(int)
	condition := expandSyntheticsCondition(d)

	log.Printf("[INFO] Creating New Relic Synthetics alert condition %s", condition.Name)

	condition, err := client.Alerts.CreateSyntheticsCondition(policyID, *condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{policyID, condition.ID}))

	return resourceNewRelicSyntheticsAlertConditionRead(d, meta)
}

func resourceNewRelicSyntheticsAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics alert condition %s", d.Id())

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	_, err = client.Alerts.GetPolicy(policyID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	condition, err := client.Alerts.GetSyntheticsCondition(policyID, id)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenSyntheticsCondition(condition, d)
}

func resourceNewRelicSyntheticsAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition := expandSyntheticsCondition(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]

	condition.ID = id

	log.Printf("[INFO] Updating New Relic Synthetics alert condition %d", id)

	_, err = client.Alerts.UpdateSyntheticsCondition(*condition)
	if err != nil {
		return err
	}

	return resourceNewRelicSyntheticsAlertConditionRead(d, meta)
}

func resourceNewRelicSyntheticsAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]

	log.Printf("[INFO] Deleting New Relic Synthetics alert condition %d", id)

	_, err = client.Alerts.DeleteSyntheticsCondition(id)

	if err != nil {
		return err
	}

	return nil
}
