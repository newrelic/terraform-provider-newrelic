package newrelic

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicSyntheticsAlertCondition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsAlertConditionCreate,
		ReadContext:   resourceNewRelicSyntheticsAlertConditionRead,
		UpdateContext: resourceNewRelicSyntheticsAlertConditionUpdate,
		DeleteContext: resourceNewRelicSyntheticsAlertConditionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "The title of this condition.",
			},
			"monitor_id": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The ID of the Synthetics monitor to be referenced in the alert condition.",
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
				Description: "Set whether to enable the alert condition. Defaults to true.",
			},
		},
	}
}

func expandSyntheticsCondition(d *schema.ResourceData) *alerts.SyntheticsCondition {
	condition := alerts.SyntheticsCondition{
		Name:      d.Get("name").(string),
		Enabled:   d.Get("enabled").(bool),
		MonitorID: d.Get("monitor_id").(string),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	return &condition
}

func flattenSyntheticsCondition(condition *alerts.SyntheticsCondition, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	_ = d.Set("policy_id", policyID)
	_ = d.Set("monitor_id", condition.MonitorID)
	_ = d.Set("name", condition.Name)
	_ = d.Set("runbook_url", condition.RunbookURL)
	_ = d.Set("enabled", condition.Enabled)

	return nil
}

func resourceNewRelicSyntheticsAlertConditionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	policyID := d.Get("policy_id").(int)
	condition := expandSyntheticsCondition(d)

	log.Printf("[INFO] Creating New Relic Synthetics alert condition %s", condition.Name)

	condition, err := client.Alerts.CreateSyntheticsCondition(policyID, *condition)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serializeIDs([]int{policyID, condition.ID}))

	return resourceNewRelicSyntheticsAlertConditionRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsAlertConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic Synthetics alert condition %s", d.Id())

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := ids[0]
	id := ids[1]

	_, err = client.Alerts.QueryPolicy(accountID, strconv.Itoa(policyID))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	condition, err := client.Alerts.GetSyntheticsCondition(policyID, id)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenSyntheticsCondition(condition, d))
}

func resourceNewRelicSyntheticsAlertConditionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	condition := expandSyntheticsCondition(d)

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return diag.FromErr(err)
	}

	id := ids[1]

	condition.ID = id

	log.Printf("[INFO] Updating New Relic Synthetics alert condition %d", id)

	_, err = client.Alerts.UpdateSyntheticsCondition(*condition)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicSyntheticsAlertConditionRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsAlertConditionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return diag.FromErr(err)
	}

	id := ids[1]

	log.Printf("[INFO] Deleting New Relic Synthetics alert condition %d", id)

	_, err = client.Alerts.DeleteSyntheticsCondition(id)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
