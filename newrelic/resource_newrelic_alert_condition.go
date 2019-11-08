package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

var alertConditionTypes = map[string][]string{
	"apm_app_metric": {
		"apdex",
		"error_percentage",
		"response_time_background",
		"response_time_web",
		"throughput_background",
		"throughput_web",
		"user_defined",
	},
	"apm_jvm_metric": {
		"cpu_utilization_time",
		"deadlocked_threads",
		"gc_cpu_time",
		"heap_memory_usage",
	},
	"apm_kt_metric": {
		"apdex",
		"error_count",
		"error_percentage",
		"response_time",
		"throughput",
	},
	"browser_metric": {
		"ajax_response_time",
		"ajax_throughput",
		"dom_processing",
		"end_user_apdex",
		"network",
		"page_rendering",
		"page_view_throughput",
		"page_views_with_js_errors",
		"request_queuing",
		"total_page_load",
		"user_defined",
		"web_application",
	},
	"mobile_metric": {
		"database",
		"images",
		"json",
		"mobile_crash_rate",
		"network_error_percentage",
		"network",
		"status_error_percentage",
		"user_defined",
		"view_loading",
	},
	"servers_metric": {
		"cpu_percentage",
		"disk_io_percentage",
		"fullest_disk_percentage",
		"load_average_one_minute",
		"memory_percentage",
		"user_defined",
	},
}

func resourceNewRelicAlertCondition() *schema.Resource {
	validAlertConditionTypes := make([]string, 0, len(alertConditionTypes))
	for k := range alertConditionTypes {
		validAlertConditionTypes = append(validAlertConditionTypes, k)
	}

	return &schema.Resource{
		Create: resourceNewRelicAlertConditionCreate,
		Read:   resourceNewRelicAlertConditionRead,
		Update: resourceNewRelicAlertConditionUpdate,
		Delete: resourceNewRelicAlertConditionDelete,
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
				Default:  true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(validAlertConditionTypes, false),
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
				//TODO: ValidateFunc from map
			},
			"runbook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"condition_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"violation_close_timer": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: intInSlice([]int{1, 2, 4, 8, 12, 24}),
			},
			"gc_metric": {
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
							ValidateFunc: intInSlice([]int{5, 10, 15, 30, 60, 120}),
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
			"user_defined_metric": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_defined_value_function": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"average", "min", "max", "total", "sample_size"}, false),
			},
		},
	}
}

func buildAlertConditionStruct(d *schema.ResourceData) *newrelic.AlertCondition {
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

	condition := newrelic.AlertCondition{
		Type:                d.Get("type").(string),
		Name:                d.Get("name").(string),
		Enabled:             d.Get("enabled").(bool),
		Entities:            entities,
		Metric:              d.Get("metric").(string),
		Terms:               terms,
		PolicyID:            d.Get("policy_id").(int),
		Scope:               d.Get("condition_scope").(string),
		ViolationCloseTimer: d.Get("violation_close_timer").(int),
		GCMetric:            d.Get("gc_metric").(string),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	if attrM, ok := d.GetOk("user_defined_metric"); ok {
		if attrVF, ok := d.GetOk("user_defined_value_function"); ok {
			condition.UserDefined = newrelic.AlertConditionUserDefined{
				Metric:        attrM.(string),
				ValueFunction: attrVF.(string),
			}
		}
	}

	return &condition
}

func readAlertConditionStruct(condition *newrelic.AlertCondition, d *schema.ResourceData) error {
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
	d.Set("type", condition.Type)
	d.Set("metric", condition.Metric)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("condition_scope", condition.Scope)
	d.Set("violation_close_timer", condition.ViolationCloseTimer)
	d.Set("gc_metric", condition.GCMetric)
	d.Set("user_defined_metric", condition.UserDefined.Metric)
	d.Set("user_defined_value_function", condition.UserDefined.ValueFunction)
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

func resourceNewRelicAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	condition := buildAlertConditionStruct(d)

	log.Printf("[INFO] Creating New Relic alert condition %s", condition.Name)

	condition, err := client.CreateAlertCondition(*condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{condition.PolicyID, condition.ID}))

	return nil
}

func resourceNewRelicAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
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

	condition, err := client.GetAlertCondition(policyID, id)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	return readAlertConditionStruct(condition, d)
}

func resourceNewRelicAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	condition := buildAlertConditionStruct(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	condition.PolicyID = policyID
	condition.ID = id

	log.Printf("[INFO] Updating New Relic alert condition %d", id)

	updatedCondition, err := client.UpdateAlertCondition(*condition)
	if err != nil {
		return err
	}

	return readAlertConditionStruct(updatedCondition, d)
}

func resourceNewRelicAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	log.Printf("[INFO] Deleting New Relic alert condition %d", id)

	if err := client.DeleteAlertCondition(policyID, id); err != nil {
		return err
	}

	return nil
}
