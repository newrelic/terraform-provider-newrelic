package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
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
