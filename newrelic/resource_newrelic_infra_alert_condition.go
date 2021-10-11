package newrelic

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				ValidateFunc: validation.StringInSlice([]string{"any", "all"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
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
		CreateContext: resourceNewRelicInfraAlertConditionCreate,
		ReadContext:   resourceNewRelicInfraAlertConditionRead,
		UpdateContext: resourceNewRelicInfraAlertConditionUpdate,
		DeleteContext: resourceNewRelicInfraAlertConditionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the alert policy where this condition should be used.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Infrastructure alert condition's name.",
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
				Description: "Whether the condition is turned on or off. Valid values are true and false. Defaults to true.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validThresholdConditionTypes, true),
				Description:  "The type of Infrastructure alert condition. Valid values are infra_process_running, infra_metric, and infra_host_not_reporting.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"event": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The metric event; for example, SystemSample or StorageSample. Supported by the infra_metric condition type.",
			},
			"where": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If applicable, this identifies any Infrastructure host filters used; for example: hostname LIKE '%cassandra%'.",
			},
			"process_where": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any filters applied to processes; for example: commandName = 'java'. Supported by the infra_process_running condition type.",
			},
			"comparison": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"above", "below", "equal"}, true),
				Description:  "The operator used to evaluate the threshold value. Valid values are above, below, and equal. Supported by the infra_metric and infra_process_running condition types.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"select": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The attribute name to identify the metric being targeted; for example, cpuPercent, diskFreePercent, or memoryResidentSizeBytes. The underlying API will automatically populate this value for Infrastructure integrations (for example diskFreePercent), so make sure to explicitly include this value to avoid diff issues. Supported by the infra_metric condition type.",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timestamp the alert condition was created.",
			},
			"updated_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timestamp the alert condition was last updated.",
			},
			"critical": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Elem:        thresholdSchema(),
				Description: "Identifies the threshold parameters for opening a critical alert violation.",
				//TODO: ValidateFunc from thresholdConditionTypes map
			},
			"warning": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				MinItems:    1,
				ForceNew:    true,
				Elem:        thresholdSchema(),
				Description: "Identifies the threshold parameters for opening a warning alert violation.",
				//TODO: ValidateFunc from thresholdConditionTypes map
			},
			"integration_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "For alerts on integrations, use this instead of event. Supported by the infra_metric condition type.",
			},
			"violation_close_timer": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      24,
				ValidateFunc: validateViolationCloseTimer(),
				Description:  "Determines how much time, in hours, will pass before a violation is automatically closed. Valid values are 1, 2, 4, 8, 12, 24, 48, or 72",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Infrastructure alert condition.",
			},
		},
	}
}

func resourceNewRelicInfraAlertConditionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	condition, err := expandInfraAlertCondition(d)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic Infra alert condition %s", condition.Name)

	condition, err = client.Alerts.CreateInfrastructureConditionWithContext(ctx, *condition)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serializeIDs([]int{condition.PolicyID, condition.ID}))

	return resourceNewRelicInfraAlertConditionRead(ctx, d, meta)
}

func resourceNewRelicInfraAlertConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := meta.(*ProviderConfig).NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic Infra alert condition %s", d.Id())

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := ids[0]
	id := ids[1]

	_, err = client.Alerts.QueryPolicyWithContext(ctx, accountID, strconv.Itoa(policyID))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	condition, err := client.Alerts.GetInfrastructureConditionWithContext(ctx, id)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenInfraAlertCondition(condition, d))
}

func resourceNewRelicInfraAlertConditionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	condition, err := expandInfraAlertCondition(d)

	if err != nil {
		return diag.FromErr(err)
	}

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := ids[0]
	id := ids[1]

	condition.PolicyID = policyID
	condition.ID = id

	log.Printf("[INFO] Updating New Relic Infra alert condition %d", id)

	_, err = client.Alerts.UpdateInfrastructureConditionWithContext(ctx, *condition)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicInfraAlertConditionRead(ctx, d, meta)
}

func resourceNewRelicInfraAlertConditionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return diag.FromErr(err)
	}

	id := ids[1]

	log.Printf("[INFO] Deleting New Relic Infra alert condition %d", id)

	if err := client.Alerts.DeleteInfrastructureConditionWithContext(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
