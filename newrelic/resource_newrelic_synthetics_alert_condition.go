package newrelic

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
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

func buildSyntheticsAlertConditionStruct(d *schema.ResourceData) *newrelic.AlertSyntheticsCondition {

	condition := newrelic.AlertSyntheticsCondition{
		Name:      d.Get("name").(string),
		Enabled:   d.Get("enabled").(bool),
		PolicyID:  d.Get("policy_id").(int),
		MonitorID: d.Get("monitor_id").(string),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	return &condition
}

func readSyntheticsAlertConditionStruct(condition *newrelic.AlertSyntheticsCondition, d *schema.ResourceData) error {
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
	client := meta.(*ProviderConfig).Client
	condition := buildSyntheticsAlertConditionStruct(d)

	log.Printf("[INFO] Creating New Relic Synthetics alert condition %s", condition.Name)

	condition, err := client.CreateAlertSyntheticsCondition(*condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{condition.PolicyID, condition.ID}))

	return resourceNewRelicSyntheticsAlertConditionRead(d, meta)
}

func resourceNewRelicSyntheticsAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic Synthetics alert condition %s", d.Id())

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	_, err = client.GetAlertPolicy(policyID)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	condition, err := client.GetAlertSyntheticsCondition(policyID, id)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	return readSyntheticsAlertConditionStruct(condition, d)
}

func resourceNewRelicSyntheticsAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	condition := buildSyntheticsAlertConditionStruct(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	condition.PolicyID = policyID
	condition.ID = id

	log.Printf("[INFO] Updating New Relic Synthetics alert condition %d", id)

	_, err = client.UpdateAlertSyntheticsCondition(*condition)
	if err != nil {
		return err
	}

	return resourceNewRelicSyntheticsAlertConditionRead(d, meta)
}

func resourceNewRelicSyntheticsAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	log.Printf("[INFO] Deleting New Relic Synthetics alert condition %d", id)

	if err := client.DeleteAlertSyntheticsCondition(policyID, id); err != nil {
		return err
	}

	return nil
}
