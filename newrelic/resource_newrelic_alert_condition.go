package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the policy where this condition should be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The title of the condition. Must be between 1 and 64 characters, inclusive.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The title of the condition. Must be between 1 and 64 characters, inclusive.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(validAlertConditionTypes, false),
				Description:  fmt.Sprintf("The type of condition. One of: (%s).", strings.Join(validAlertConditionTypes, ", ")),
			},
			"entities": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
				MinItems:    1,
				Description: "The instance IDs associated with this condition.",
			},
			"metric": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The metric field accepts parameters based on the type set.",
				//TODO: ValidateFunc from map
			},
			"runbook_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Runbook URL to display in notifications.",
			},
			"condition_scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"application", "instance"}, false),
				Description:  "One of (application, instance). Choose application for most scenarios. If you are using the JVM plugin in New Relic, the instance setting allows your condition to trigger for specific app instances.",
			},
			"violation_close_timer": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: intInSlice([]int{1, 2, 4, 8, 12, 24}),
				Description:  "Automatically close instance-based violations, including JVM health metric violations, after the number of hours specified. Must be: 1, 2, 4, 8, 12 or 24.",
			},
			"gc_metric": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A valid Garbage Collection metric e.g. GC/G1 Young Generation. This is required if you are using apm_jvm_metric with gc_cpu_time condition type.",
			},
			"term": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(5, 120),
							Description:  "In minutes, must be in the range of 5 to 120, inclusive.",
						},
						"operator": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "equal",
							ValidateFunc: validation.StringInSlice([]string{"above", "below", "equal"}, false),
							Description:  "One of (above, below, equal). Defaults to equal.",
						},
						"priority": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "critical",
							ValidateFunc: validation.StringInSlice([]string{"critical", "warning"}, false),
							Description:  "One of (critical, warning). Defaults to critical.",
						},
						"threshold": {
							Type:         schema.TypeFloat,
							Required:     true,
							ValidateFunc: float64Gte(0.0),
							Description:  "Must be 0 or greater.",
						},
						"time_function": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
							Description:  "One of (all, any).",
						},
					},
				},
				Required: true,
				MinItems: 1,
			},
			"user_defined_metric": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A custom metric to be evaluated.",
			},
			"user_defined_value_function": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"average", "min", "max", "total", "sample_size"}, false),
				Description:  "One of: (average, min, max, total, sample_size).",
			},
		},
	}
}

func resourceNewRelicAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition, err := expandAlertCondition(d)
	if err != nil {
		return err
	}

	policyID := d.Get("policy_id").(int)

	log.Printf("[INFO] Creating New Relic alert condition %s", condition.Name)

	condition, err = client.Alerts.CreateCondition(policyID, *condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{policyID, condition.ID}))

	return nil
}

func resourceNewRelicAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)
	client := meta.(*ProviderConfig).NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic alert condition %s", d.Id())

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	_, err = client.Alerts.QueryPolicy(accountID, strconv.Itoa(policyID))
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

	d.Set("policy_id", policyID)

	return flattenAlertCondition(condition, d)
}

func resourceNewRelicAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition, err := expandAlertCondition(d)
	if err != nil {
		return err
	}

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]
	condition.ID = id

	log.Printf("[INFO] Updating New Relic alert condition %d", id)

	updatedCondition, err := client.Alerts.UpdateCondition(*condition)
	if err != nil {
		return err
	}

	d.Set("policy_id", policyID)

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
