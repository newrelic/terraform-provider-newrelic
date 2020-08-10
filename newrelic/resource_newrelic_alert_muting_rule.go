package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicAlertMutingRule() *schema.Resource {
	// validAlertConditionTypes := make([]string, 0, len(alertConditionTypes))
	// for k := range alertConditionTypes {
	// 	validAlertConditionTypes = append(validAlertConditionTypes, k)
	// }

	return &schema.Resource{
		Create: resourceNewRelicAlertMutingRuleCreate,
		Read:   resourceNewRelicAlertMutingRuleRead,
		Update: resourceNewRelicAlertMutingRuleUpdate,
		Delete: resourceNewRelicAlertMutingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"accountId": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The muting rule's account Id.",
			},
			"condition": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The condition that defines which violations to target.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The individual MutingRuleConditions within the group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The attribute on a violation.",
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
				Description: "Whether the MutingRule is enabled",
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
		},
	}
}

func expandMutingRuleCondition(cfg []interface{}) []alerts.MutingRuleConditionGroup {
	if len(cfg) == 0 {
		return []alerts.MutingRuleConditionGroup{}
	}

	conditionGroup := []alerts.MutingRuleConditionGroup{}

	conditionGroup.Conditions = expandMutingRuleConditions(cfg)

	if operator, ok := cfg["operator"]; ok {
		conditionGroup.Operator = operator.(string)
	}

	return conditionGroup
}

func expandMutingRuleConditions(cfg []interface{}) []alerts.MutingRuleCondition {

	conditions := []alerts.MutingRuleCondition{}

	for i, rawCfg := range cfg {
		cfg := rawCfg.(map[string]interface{})
		if attribute, ok := cfg["attribute"]; ok {
			conditions.Attribute = attribute.(string)
		}

		if operator, ok := cfg["operator"]; ok {
			conditions.Operator = operator.(string)
		}
		if values, ok := cfg["values"]; ok {
			conditions.Values = values.(list)
		}
	}

	return conditions
}

func expandMutingRuleCreateInput(d *schema.ResourceData) alerts.MutingRuleCreateInput {
	createInput := alerts.MutingRuleCreateInput{
		// Account:     d.Get("accountId").(int),
		Enabled:     d.Get("enabled").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if e, ok := d.GetOk("Condition"); ok {
		createInput.Condition = expandMutingRuleCondition(e.(*schema.Set).List())
	}

	return createInput
}

func resourceNewRelicAlertMutingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient // Getting instance of the client
	expanded := expandMutingRuleCreateInput(d)

	accountID := d.Get("account_id").(int) // Account ID from the schema

	log.Printf("[INFO] Creating New Relic MutingRule alerts")

	created, err := client.Alerts.CreateMutingRule(accountID, expanded)
	if err != nil {
		return err
	}

	ids := mutingRuleIDs{
		AccountID: accountID,
		ID:        created.ID,
	}

	// Add helper function here.
	d.SetId(ids.String())

	return nil
}

func resourceNewRelicAlertMutingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic MutingRule alerts")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	mutingRule, err := client.Alerts.GetMutingRule(ids.accountId, ids.rule)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenAlertMutingRule(mutingRule, d)

}

// func resourceNewRelicAlertMutingRuleUpdate(d *schema.ResourceData, meta interface{}) error {}

// func resourceNewRelicAlertMutingRuleDelete(d *schema.ResourceData, meta interface{}) error {}

type mutingRuleIDs struct {
	AccountID int
	ID        int
}
