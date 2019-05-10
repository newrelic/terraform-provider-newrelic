package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func resourceNewRelicNrqlAlertCondition() *schema.Resource {

	return &schema.Resource{
		Create: resourceNewRelicNrqlAlertConditionCreate,
		Read:   resourceNewRelicNrqlAlertConditionRead,
		Update: resourceNewRelicNrqlAlertConditionUpdate,
		Delete: resourceNewRelicNrqlAlertConditionDelete,
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
			"expected_groups": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ignore_overlap": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nrql": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:     schema.TypeString,
							Required: true,
						},
						"since_value": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								valueString := val.(string)
								v, err := strconv.Atoi(valueString)
								if err != nil {
									errs = append(errs, fmt.Errorf("Error converting string to int: %#v", err))
								}
								if v < 1 || v > 20 {
									errs = append(errs, fmt.Errorf("%q must be between 0 and 20 inclusive, got: %d", key, v))
								}
								return
							},
						},
					},
				},
			},
			"term": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"operator": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "equal",
							ValidateFunc: validation.StringInSlice([]string{"above", "below", "equal"}, false),
						},
						"priority": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "critical",
							ValidateFunc: validation.StringInSlice([]string{"critical", "warning"}, false),
						},
						"threshold": {
							Type:         schema.TypeFloat,
							Required:     true,
							ValidateFunc: float64Gte(0.0),
						},
						"time_function": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
						},
					},
				},
				Required: true,
				MinItems: 1,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "static",
				ValidateFunc: validation.StringInSlice([]string{"static", "outlier", "baseline"}, false),
			},
			"value_function": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "single_value",
				ValidateFunc: validation.StringInSlice([]string{"single_value", "sum"}, false),
			},
		},
	}
}

func buildNrqlAlertConditionStruct(d *schema.ResourceData) *newrelic.AlertNrqlCondition {
	termSet := d.Get("term").([]interface{})
	terms := make([]newrelic.AlertConditionTerm, len(termSet))

	for i, termI := range termSet {
		termM := termI.(map[string]interface{})

		terms[i] = newrelic.AlertConditionTerm{
			Duration:     termM["duration"].(int),
			Operator:     termM["operator"].(string),
			Priority:     termM["priority"].(string),
			Threshold:    termM["threshold"].(float64),
			TimeFunction: termM["time_function"].(string),
		}
	}

	query := newrelic.AlertNrqlQuery{}

	if nrqlQuery, ok := d.GetOk("nrql.0.query"); ok {
		query.Query = nrqlQuery.(string)
	}

	if sinceValue, ok := d.GetOk("nrql.0.since_value"); ok {
		query.SinceValue = sinceValue.(string)
	}

	condition := newrelic.AlertNrqlCondition{
		Name:          d.Get("name").(string),
		Type:          d.Get("type").(string),
		Enabled:       d.Get("enabled").(bool),
		Terms:         terms,
		PolicyID:      d.Get("policy_id").(int),
		Nrql:          query,
		ValueFunction: d.Get("value_function").(string),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	if attr, ok := d.GetOkExists("ignore_overlap"); ok {
		condition.IgnoreOverlap = attr.(bool)
	}

	if attr, ok := d.GetOk("expected_groups"); ok {
		condition.ExpectedGroups = attr.(int)
	}

	return &condition
}

func readNrqlAlertConditionStruct(condition *newrelic.AlertNrqlCondition, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	d.Set("policy_id", policyID)
	d.Set("name", condition.Name)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("enabled", condition.Enabled)
	d.Set("nrql.0.Query", condition.Nrql.Query)
	d.Set("nrql.0.SinceValue", condition.Nrql.SinceValue)

	if condition.ValueFunction == "" {
		d.Set("value_function", "single_value")
	} else {
		d.Set("value_function", condition.ValueFunction)
	}

	var terms []map[string]interface{}

	for _, src := range condition.Terms {
		dst := map[string]interface{}{
			"duration":      src.Duration,
			"operator":      src.Operator,
			"priority":      src.Priority,
			"threshold":     src.Threshold,
			"time_function": src.TimeFunction,
		}
		terms = append(terms, dst)
	}

	if err := d.Set("term", terms); err != nil {
		return fmt.Errorf("[DEBUG] Error setting alert condition terms: %#v", err)
	}

	return nil
}

func resourceNewRelicNrqlAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	condition := buildNrqlAlertConditionStruct(d)

	log.Printf("[INFO] Creating New Relic NRQL alert condition %s", condition.Name)

	condition, err := client.CreateAlertNrqlCondition(*condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{condition.PolicyID, condition.ID}))

	return resourceNewRelicNrqlAlertConditionRead(d, meta)
}

func resourceNewRelicNrqlAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic NRQL alert condition %s", d.Id())

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	condition, err := client.GetAlertNrqlCondition(policyID, id)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	return readNrqlAlertConditionStruct(condition, d)
}

func resourceNewRelicNrqlAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	condition := buildNrqlAlertConditionStruct(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	condition.PolicyID = policyID
	condition.ID = id

	log.Printf("[INFO] Updating New Relic NRQL alert condition %d", id)

	_, err = client.UpdateAlertNrqlCondition(*condition)
	if err != nil {
		return err
	}

	return resourceNewRelicNrqlAlertConditionRead(d, meta)
}

func resourceNewRelicNrqlAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	log.Printf("[INFO] Deleting New Relic NRQL alert condition %d", id)

	if err := client.DeleteAlertNrqlCondition(policyID, id); err != nil {
		return err
	}

	return nil
}
