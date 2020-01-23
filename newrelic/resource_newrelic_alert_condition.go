package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
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

func resourceNewRelicAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition := expandAlertCondition(d)

	log.Printf("[INFO] Creating New Relic alert condition %s", condition.Name)

	condition, err := client.Alerts.CreateCondition(*condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{condition.PolicyID, condition.ID}))

	return nil
}

func resourceNewRelicAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic alert condition %s", d.Id())

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

	condition, err := client.Alerts.GetCondition(policyID, id)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenAlertCondition(condition, d)
}

func resourceNewRelicAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition := expandAlertCondition(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	condition.PolicyID = policyID
	condition.ID = id

	log.Printf("[INFO] Updating New Relic alert condition %d", id)

	updatedCondition, err := client.Alerts.UpdateCondition(*condition)
	if err != nil {
		return err
	}

	return flattenAlertCondition(updatedCondition, d)
}

func resourceNewRelicAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]

	log.Printf("[INFO] Deleting New Relic alert condition %d", id)

	_, err = client.Alerts.DeleteCondition(id)
	if err != nil {
		return err
	}

	return nil
}
