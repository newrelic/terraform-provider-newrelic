package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

var thresholdConditionTypes = map[string][]string{
	"infra_process_running": {
		"duration",
		"value",
	},
	"infra_metric": {
		"duration",
		"value",
		"time_function",
	},
	"infra_host_not_reporting": {
		"duration",
	},
}

// thresholdSchema returns the schema to use for threshold.
func thresholdSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
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
				Computed: true,
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
			"violation_close_timer": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      24,
				ValidateFunc: intInSlice([]int{0, 1, 2, 4, 8, 12, 24, 48, 72}),
			},
		},
	}
}

func resourceNewRelicInfraAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition, err := expandInfraAlertCondition(d)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating New Relic Infra alert condition %s", condition.Name)

	condition, err = client.Alerts.CreateInfrastructureCondition(*condition)

	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{condition.PolicyID, condition.ID}))

	return resourceNewRelicInfraAlertConditionRead(d, meta)
}

func resourceNewRelicInfraAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Infra alert condition %s", d.Id())

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

	condition, err := client.Alerts.GetInfrastructureCondition(id)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenInfraAlertCondition(condition, d)
}

func resourceNewRelicInfraAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition, err := expandInfraAlertCondition(d)

	if err != nil {
		return err
	}

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	condition.PolicyID = policyID
	condition.ID = id

	log.Printf("[INFO] Updating New Relic Infra alert condition %d", id)

	_, err = client.Alerts.UpdateInfrastructureCondition(*condition)
	if err != nil {
		return err
	}

	return resourceNewRelicInfraAlertConditionRead(d, meta)
}

func resourceNewRelicInfraAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]

	log.Printf("[INFO] Deleting New Relic Infra alert condition %d", id)

	if err := client.Alerts.DeleteInfrastructureCondition(id); err != nil {
		return err
	}

	return nil
}
