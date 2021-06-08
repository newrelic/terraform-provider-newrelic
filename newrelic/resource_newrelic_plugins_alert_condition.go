package newrelic

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNewRelicPluginsAlertCondition() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "`newrelic_plugins_alert_condition` has been deprecated and will not be supported as of June 16, 2021.  Use `newrelic_nrql_alert_condition` instead.",
		Create:             resourceNewRelicPluginsAlertConditionCreate,
		Read:               resourceNewRelicPluginsAlertConditionRead,
		Update:             resourceNewRelicPluginsAlertConditionUpdate,
		Delete:             resourceNewRelicPluginsAlertConditionDelete,
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
				Description: "Whether or not this condition is enabled.",
			},
			"entities": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
				MinItems:    1,
				Description: "The plugin component IDs to target.",
			},
			"metric": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plugin metric to evaluate.",
			},
			"metric_description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The metric description.",
			},
			"value_function": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"min", "max", "average", "sample_size", "total", "percent"}, false),
				Description:  "The value function to apply to the metric data.  One of `min`, `max`, `average`, `sample_size`, `total`, or `percent`.",
			},
			"runbook_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Runbook URL to display in notifications.",
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
							Description:  "One of `above`, `below`, or `equal`. Defaults to equal.",
						},
						"priority": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "critical",
							ValidateFunc: validation.StringInSlice([]string{"critical", "warning"}, false),
							Description:  "One of `critical` or `warning`. Defaults to critical.",
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
							Description:  "One of `all` or `any`.",
						},
					},
				},
				Required: true,
				MinItems: 1,
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the installed plugin instance which produces the metric.",
			},
			"plugin_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the plugin which produces the metric.",
			},
		},
	}
}

func resourceNewRelicPluginsAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("plugins have reached end of life, use `newrelic_nrql_alert_condition` instead")
}

func resourceNewRelicPluginsAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	return errors.New("plugins have reached end of life, use `newrelic_nrql_alert_condition` instead")
}

func resourceNewRelicPluginsAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("plugins have reached end of life, use `newrelic_nrql_alert_condition` instead")
}

func resourceNewRelicPluginsAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	return errors.New("plugins have reached end of life, use `newrelic_nrql_alert_condition` instead")
}
