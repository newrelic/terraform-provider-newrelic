package newrelic

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

var thresholdConditionTypes = map[string][]string{
	"infra_process_running": {
		"duration_minutes",
		"value",
	},
	"infra_metric": {
		"duration_minutes",
		"value",
		"time_function",
	},
	"infra_host_not_reporting": {
		"duration_minutes",
	},
}

// thresholdSchema returns the schema to use for threshold.
//
func thresholdSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"duration": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 60),
			},
			"time_function": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"any", "all"}, false),
			},
		},
	}
}

func resourceNewRelicInfraAlertCondition() *schema.Resource {
	validThresholdConditionTypes := make([]string, 0, len(thresholdConditionTypes))
	for k := range thresholdConditionTypes {
		validThresholdConditionTypes = append(validThresholdConditionTypes, k)
	}

	return &schema.Resource{
		Create: resourceNewRelicInfraAlertConditionCreate,
		Read:   resourceNewRelicInfraAlertConditionRead,
		Update: resourceNewRelicInfraAlertConditionUpdate,
		Delete: resourceNewRelicInfraAlertConditionDelete,
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
			"runbook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validThresholdConditionTypes, false),
			},
			"event": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"where": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"process_where": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comparison": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"above", "below", "equal"}, false),
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"critical": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				Elem:     thresholdSchema(),
				//TODO: ValidateFunc from thresholdConditionTypes map
			},
			"warning": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				MinItems: 1,
				ForceNew: true,
				Elem:     thresholdSchema(),
				//TODO: ValidateFunc from thresholdConditionTypes map
			},
			"integration_provider": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildInfraAlertConditionStruct(d *schema.ResourceData) *newrelic.AlertInfraCondition {

	condition := newrelic.AlertInfraCondition{
		Name:       d.Get("name").(string),
		Enabled:    d.Get("enabled").(bool),
		PolicyID:   d.Get("policy_id").(int),
		Event:      d.Get("event").(string),
		Comparison: d.Get("comparison").(string),
		Select:     d.Get("select").(string),
		Type:       d.Get("type").(string),
		Critical:   expandAlertThreshold(d.Get("critical")),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}
	if attr, ok := d.GetOk("warning"); ok {
		condition.Warning = expandAlertThreshold(attr)
	}

	if attr, ok := d.GetOk("where"); ok {
		condition.Where = attr.(string)
	}

	if attr, ok := d.GetOk("process_where"); ok {
		condition.ProcessWhere = attr.(string)
	}

	if attr, ok := d.GetOk("integration_provider"); ok {
		condition.IntegrationProvider = attr.(string)
	}

	return &condition
}

func readInfraAlertConditionStruct(condition *newrelic.AlertInfraCondition, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	d.Set("policy_id", policyID)
	d.Set("name", condition.Name)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("enabled", condition.Enabled)
	d.Set("created_at", condition.CreatedAt)
	d.Set("updated_at", condition.UpdatedAt)

	if condition.Where != "" {
		d.Set("where", condition.Where)
	}

	if condition.ProcessWhere != "" {
		d.Set("process_where", condition.ProcessWhere)
	}

	if condition.IntegrationProvider != "" {
		d.Set("integration_provider", condition.IntegrationProvider)
	}

	if err := d.Set("critical", flattenAlertThreshold(condition.Critical)); err != nil {
		return err
	}

	if condition.Warning != nil {
		if err := d.Set("warning", flattenAlertThreshold(condition.Warning)); err != nil {
			return err
		}
	}

	return nil
}

func resourceNewRelicInfraAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).InfraClient
	condition := buildInfraAlertConditionStruct(d)

	log.Printf("[INFO] Creating New Relic Infra alert condition %s", condition.Name)

	condition, err := client.CreateAlertInfraCondition(*condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{condition.PolicyID, condition.ID}))

	return resourceNewRelicInfraAlertConditionRead(d, meta)
}

func resourceNewRelicInfraAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).InfraClient
	policyClient := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic Infra alert condition %s", d.Id())

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	_, err = policyClient.GetAlertPolicy(policyID)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	condition, err := client.GetAlertInfraCondition(policyID, id)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	return readInfraAlertConditionStruct(condition, d)
}

func resourceNewRelicInfraAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).InfraClient
	condition := buildInfraAlertConditionStruct(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	condition.PolicyID = policyID
	condition.ID = id

	log.Printf("[INFO] Updating New Relic Infra alert condition %d", id)

	_, err = client.UpdateAlertInfraCondition(*condition)
	if err != nil {
		return err
	}

	return resourceNewRelicInfraAlertConditionRead(d, meta)
}

func resourceNewRelicInfraAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).InfraClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	log.Printf("[INFO] Deleting New Relic Infra alert condition %d", id)

	if err := client.DeleteAlertInfraCondition(policyID, id); err != nil {
		return err
	}

	return nil
}

func expandAlertThreshold(v interface{}) *newrelic.AlertInfraThreshold {
	rah := v.([]interface{})[0].(map[string]interface{})

	alertInfraThreshold := &newrelic.AlertInfraThreshold{
		Duration: rah["duration"].(int),
	}

	if val, ok := rah["value"]; ok {
		alertInfraThreshold.Value = val.(int)
	}

	if val, ok := rah["time_function"]; ok {
		alertInfraThreshold.Function = val.(string)
	}

	return alertInfraThreshold
}

func flattenAlertThreshold(v *newrelic.AlertInfraThreshold) []interface{} {
	alertInfraThreshold := map[string]interface{}{
		"duration":      v.Duration,
		"value":         v.Value,
		"time_function": v.Function,
	}

	return []interface{}{alertInfraThreshold}
}
