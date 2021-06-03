package newrelic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
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
	client := meta.(*ProviderConfig).NewClient
	condition := expandPluginsCondition(d)
	policyID := d.Get("policy_id").(int)

	log.Printf("[INFO] Creating New Relic alert condition %s", condition.Name)

	condition, err := client.Alerts.CreatePluginsCondition(policyID, *condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{policyID, condition.ID}))

	return nil
}

func resourceNewRelicPluginsAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
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

	condition, err := client.Alerts.GetPluginsCondition(policyID, id)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("policy_id", policyID)

	return flattenPluginsCondition(condition, d)
}

func resourceNewRelicPluginsAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition := expandPluginsCondition(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]
	condition.ID = id

	log.Printf("[INFO] Updating New Relic alert condition %d", id)

	updatedCondition, err := client.Alerts.UpdatePluginsCondition(*condition)
	if err != nil {
		return err
	}

	return flattenPluginsCondition(updatedCondition, d)
}

func resourceNewRelicPluginsAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]

	log.Printf("[INFO] Deleting New Relic alert condition %d", id)

	_, err = client.Alerts.DeletePluginsCondition(id)

	if err != nil {
		return err
	}

	return nil
}
