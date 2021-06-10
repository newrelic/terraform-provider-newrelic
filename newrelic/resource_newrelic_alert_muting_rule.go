package newrelic

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func validateMutingRuleConditionAttribute(val interface{}, key string) (warns []string, errs []error) {
	valueString := val.(string)
	attemptedTagRegex := regexp.MustCompile(`^tag`)
	correctTagRegex := regexp.MustCompile(`^tag\..+$`)

	// tag.SomeValue attempted but does not match allowed format
	if attemptedTagRegex.Match([]byte(valueString)) {
		if !correctTagRegex.Match([]byte(valueString)) {
			errs = append(errs, fmt.Errorf("%#v of %#v must be in the format tag.tag_name", key, valueString))
			return
		}
		return
	}
	v := validation.StringInSlice([]string{"accountId", "conditionId", "policyId", "policyName", "conditionName", "conditionType", "conditionRunbookUrl", "product", "targetId", "targetName", "nrqlEventType", "tag", "nrqlQuery"}, false)
	return v(valueString, key)
}
func validateNaiveDateTime(val interface{}, key string) (warns []string, errs []error) {
	valueString := val.(string)

	// test conversion to desired format:
	_, err := time.Parse("2006-01-02T15:04:05", valueString)
	if err != nil {
		errs = append(errs, fmt.Errorf("%#v of %#v must be in the format 2006-01-02T15:04:05", key, valueString))
	}
	return
}

func scheduleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"end_repeat": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The datetime stamp when the MutingRule schedule should stop repeating.",
				ConflictsWith: []string{"schedule.0.repeat_count"},
				ValidateFunc:  validateNaiveDateTime,
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The datetime stamp representing when the MutingRule should end.",
				ValidateFunc: validateNaiveDateTime,
			},
			"repeat": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The frequency the MutingRule schedule repeats. One of [DAILY, WEEKLY, MONTHLY]",
				ValidateFunc: validation.StringInSlice([]string{"DAILY", "WEEKLY", "MONTHLY"}, false),
			},
			"repeat_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The number of times the MutingRule schedule should repeat.",
				ConflictsWith: []string{"schedule.0.end_repeat"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
			"start_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The datetime stamp representing when the MutingRule should start.",
				ValidateFunc: validateNaiveDateTime,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The time zone that applies to the MutingRule schedule.",
			},
			"weekly_repeat_days": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringInSlice([]string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"}, false)},
				Optional:    true,
				Description: "The day(s) of the week that a MutingRule should repeat when the repeat field is set to WEEKLY.",
				MinItems:    0,
				MaxItems:    7,
			},
		},
	}
}

func resourceNewRelicAlertMutingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAlertMutingRuleCreate,
		ReadContext:   resourceNewRelicAlertMutingRuleRead,
		UpdateContext: resourceNewRelicAlertMutingRuleUpdate,
		DeleteContext: resourceNewRelicAlertMutingRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The account id of the MutingRule..",
			},
			"condition": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The condition that defines which violations to target.",
				MaxItems:    1,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The individual MutingRuleConditions within the group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateMutingRuleConditionAttribute,
										Description:  "The attribute on a violation.",
									},
									"operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The operator used to compare the attribute's value with the supplied value(s).",
									},
									"values": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The value(s) to compare against the attribute's value.",
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The operator used to combine all the MutingRuleConditions within the group.",
						},
					},
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the MutingRule is enabled.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the MutingRule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the MutingRule.",
			},
			"schedule": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Elem:        scheduleSchema(),
				Description: "The time window when the MutingRule should actively mute violations.",
			},
		},
	}
}

func resourceNewRelicAlertMutingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	createInput, err := expandMutingRuleCreateInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Creating New Relic alert muting rule.")

	created, err := client.Alerts.CreateMutingRule(accountID, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serializeIDs([]int{accountID, created.ID}))

	return resourceNewRelicAlertMutingRuleRead(ctx, d, meta)
}

func resourceNewRelicAlertMutingRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic alert muting rule.")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := ids[0]
	mutingRuleID := ids[1]

	mutingRule, err := client.Alerts.GetMutingRule(accountID, mutingRuleID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenMutingRule(mutingRule, d))
}

func resourceNewRelicAlertMutingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Updating New Relic alert muting rule.")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	mutingRuleID := ids[1]

	updateInput, err := expandMutingRuleUpdateInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Alerts.UpdateMutingRule(accountID, mutingRuleID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicAlertMutingRuleRead(ctx, d, meta)
}

func resourceNewRelicAlertMutingRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One muting rule alert.")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := ids[0]
	mutingRuleID := ids[1]

	err = client.Alerts.DeleteMutingRule(accountID, mutingRuleID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
