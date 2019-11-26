package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func resourceNewRelicPluginsAlertCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicPluginsAlertConditionCreate,
		Read:   resourceNewRelicPluginsAlertConditionRead,
		Update: resourceNewRelicPluginsAlertConditionUpdate,
		Delete: resourceNewRelicPluginsAlertConditionDelete,
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"entities": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Required: true,
				MinItems: 1,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value_function": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"min", "max", "average", "sample_size", "total", "percent"}, false),
			},
			"runbook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"term": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(5, 120),
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
			"plugin_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plugin_guid": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildPluginsAlertConditionStruct(d *schema.ResourceData) *newrelic.AlertPluginsCondition {
	entitySet := d.Get("entities").(*schema.Set).List()
	entities := make([]string, len(entitySet))

	for i, entity := range entitySet {
		entities[i] = strconv.Itoa(entity.(int))
	}

	termSet := d.Get("term").(*schema.Set).List()
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

	plugin := newrelic.AlertPlugin{ID: d.Get("plugin_id").(string), GUID: d.Get("plugin_guid").(string)}

	condition := newrelic.AlertPluginsCondition{
		Name:              d.Get("name").(string),
		Enabled:           d.Get("enabled").(bool),
		Entities:          entities,
		Metric:            d.Get("metric").(string),
		MetricDescription: d.Get("metric_description").(string),
		ValueFunction:     d.Get("value_function").(string),
		Terms:             terms,
		PolicyID:          d.Get("policy_id").(int),
		Plugin:            plugin,
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	return &condition
}

func readPluginsAlertConditionStruct(condition *newrelic.AlertPluginsCondition, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	entities := make([]int, len(condition.Entities))
	for i, entity := range condition.Entities {
		v, err := strconv.ParseInt(entity, 10, 32)
		if err != nil {
			return err
		}
		entities[i] = int(v)
	}

	d.Set("policy_id", policyID)
	d.Set("name", condition.Name)
	d.Set("enabled", condition.Enabled)
	d.Set("metric", condition.Metric)
	d.Set("metric_description", condition.MetricDescription)
	d.Set("value_function", condition.ValueFunction)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("plugin_id", condition.Plugin.ID)
	d.Set("plugin_guid", condition.Plugin.GUID)
	if err := d.Set("entities", entities); err != nil {
		return fmt.Errorf("[DEBUG] Error setting alert condition entities: %#v", err)
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

func resourceNewRelicPluginsAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	condition := buildPluginsAlertConditionStruct(d)

	log.Printf("[INFO] Creating New Relic alert condition %s", condition.Name)

	condition, err := client.CreateAlertPluginsCondition(*condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{condition.PolicyID, condition.ID}))

	return nil
}

func resourceNewRelicPluginsAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic alert condition %s", d.Id())

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

	condition, err := client.GetAlertPluginsCondition(policyID, id)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	return readPluginsAlertConditionStruct(condition, d)
}

func resourceNewRelicPluginsAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	condition := buildPluginsAlertConditionStruct(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	condition.PolicyID = policyID
	condition.ID = id

	log.Printf("[INFO] Updating New Relic alert condition %d", id)

	updatedCondition, err := client.UpdateAlertPluginsCondition(*condition)
	if err != nil {
		return err
	}

	return readPluginsAlertConditionStruct(updatedCondition, d)
}

func resourceNewRelicPluginsAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	log.Printf("[INFO] Deleting New Relic alert condition %d", id)

	if err := client.DeleteAlertPluginsCondition(policyID, id); err != nil {
		return err
	}

	return nil
}
