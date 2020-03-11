package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
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
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the policy where this condition should be used.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the condition.",
			},
			"runbook_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Runbook URL to display in notifications.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable the alert condition.",
			},
			"expected_groups": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of expected groups when using outlier detection.",
			},
			"ignore_overlap": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to look for a convergence of groups when using outlier detection.",
			},
			"violation_time_limit_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{3600, 7200, 14400, 28800, 43200, 86400}),
				Description:  "Sets a time limit, in seconds, that will automatically force-close a long-lasting violation after the time limit you select. Possible values are 3600, 7200, 14400, 28800, 43200, and 86400.",
			},
			"nrql": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "A NRQL query.",
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
									errs = append(errs, fmt.Errorf("error converting string to int: %#v", err))
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
				Type:        schema.TypeSet,
				Description: "A list of terms for this condition. ",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 120),
							Description:  "In minutes, must be in the range of 1 to 120, inclusive.",
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
				Description:  "Possible values are single_value, sum.",
			},
		},
	}
}

func resourceNewRelicNrqlAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition := expandNrqlAlertConditionStruct(d)
	policyID := d.Get("policy_id").(int)

	log.Printf("[INFO] Creating New Relic NRQL alert condition %s", condition.Name)

	condition, err := client.Alerts.CreateNrqlCondition(policyID, *condition)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{policyID, condition.ID}))

	return resourceNewRelicNrqlAlertConditionRead(d, meta)
}

func resourceNewRelicNrqlAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic NRQL alert condition %s", d.Id())

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

	condition, err := client.Alerts.GetNrqlCondition(policyID, id)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("policy_id", policyID)

	return flattenNrqlConditionStruct(condition, d)
}

func resourceNewRelicNrqlAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition := expandNrqlAlertConditionStruct(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]
	condition.ID = id

	log.Printf("[INFO] Updating New Relic NRQL alert condition %d", id)

	_, err = client.Alerts.UpdateNrqlCondition(*condition)
	if err != nil {
		return err
	}

	return resourceNewRelicNrqlAlertConditionRead(d, meta)
}

func resourceNewRelicNrqlAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]

	log.Printf("[INFO] Deleting New Relic NRQL alert condition %d", id)

	_, err = client.Alerts.DeleteNrqlCondition(id)
	if err != nil {
		return err
	}

	return nil
}
