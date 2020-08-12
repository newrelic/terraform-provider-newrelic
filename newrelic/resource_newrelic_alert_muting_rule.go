package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicAlertMutingRule() *schema.Resource {

	return &schema.Resource{
		Create: resourceNewRelicAlertMutingRuleCreate,
		Read:   resourceNewRelicAlertMutingRuleRead,
		Update: resourceNewRelicAlertMutingRuleUpdate,
		Delete: resourceNewRelicAlertMutingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
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

func expandMutingRuleValues(values []interface{}) []string {
	perms := make([]string, len(values))

	for i, values := range values {
		perms[i] = values.(string)
	}

	return perms
}

func expandMutingRuleCondition(cfg interface{}) alerts.MutingRuleCondition {

	conditions := alerts.MutingRuleCondition{}

	conditionCfg := cfg.(map[string]interface{})
	condition := alerts.MutingRuleCondition{}

	if attribute, ok := conditionCfg["attribute"]; ok {
		condition.Attribute = attribute.(string)
	}

	if operator, ok := conditionCfg["operator"]; ok {
		condition.Operator = operator.(string)
	}
	if values, ok := conditionCfg["values"]; ok {
		condition.Values = expandMutingRuleValues(values.([]interface{}))
	}

	return conditions
}

func expandMutingRuleConditionGroup(cfg map[string]interface{}) alerts.MutingRuleConditionGroup {

	conditionGroup := alerts.MutingRuleConditionGroup{}
	var expandedConditions []alerts.MutingRuleCondition

	fmt.Println("====================")
	fmt.Printf("cfg: %+v \n", cfg)
	fmt.Println("====================")

	conditions := cfg["conditions"].(*schema.Set).List()

	for _, c := range conditions {

		var y = expandMutingRuleCondition(c)
		expandedConditions = append(expandedConditions, y)
	}

	conditionGroup.Conditions = expandedConditions

	if operator, ok := cfg["operator"]; ok {
		conditionGroup.Operator = operator.(string)
	}

	return conditionGroup
}

func expandMutingRuleCreateInput(d *schema.ResourceData) alerts.MutingRuleCreateInput {
	createInput := alerts.MutingRuleCreateInput{
		Enabled:     d.Get("enabled").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if e, ok := d.GetOk("condition"); ok {
		createInput.Condition = expandMutingRuleConditionGroup(e.([]interface{})[0].(map[string]interface{}))
	}

	return createInput
}

func resourceNewRelicAlertMutingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient // Getting instance of the client
	expanded := expandMutingRuleCreateInput(d)

	fmt.Println("====================")
	fmt.Printf("expanded: %+v \n", expanded)
	fmt.Println("====================")

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

	d.SetId(ids.String())

	return resourceNewRelicAlertMutingRuleRead(d, meta)
}

// func flattenMutingRuleValues(in []alerts.MutingRuleCondition) []string {
// 	var values []map[string]interface{}

// 	for _, src := range *in {
// 		dst := map[string]interface{}{
// 			"value": src.Value,
// 		}
// 		values = append(values, dst)
// 	}

// 	return values
// }

func flattenMutingRuleCondition(in alerts.MutingRuleCondition) []map[string]interface{} {
	var condition []map[string]interface{}

	// for _, src := range *in {
	dst := map[string]interface{}{
		"attribute": in.Attribute,
		"operator":  in.Operator,
		"values":    in.Values,
	}
	condition = append(condition, dst)
	// }

	return condition

}

func flattenMutingRuleConditionGroup(in alerts.MutingRuleConditionGroup) []map[string]interface{} {
	var condition []map[string]interface{}

	for _, src := range in.Conditions {
		dst := map[string]interface{}{
			"operator":   src.Operator,
			"conditions": flattenMutingRuleCondition(src),
		}
		condition = append(condition, dst)
	}

	return condition
}

func flattenMutingRule(mutingRule *alerts.MutingRule, d *schema.ResourceData) error {
	d.Set("enabled", mutingRule.Enabled)
	d.Set("condition", flattenMutingRuleConditionGroup(mutingRule.Condition))
	d.Set("description", mutingRule.Description)
	d.Set("name", mutingRule.Name)

	return nil
}

func resourceNewRelicAlertMutingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic MutingRule alerts")

	ids, err := parseMutingRuleIDs(d.Id())
	if err != nil {
		return err
	}

	mutingRule, err := client.Alerts.GetMutingRule(ids.AccountID, ids.ID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenMutingRule(mutingRule, d)

}

func resourceNewRelicAlertMutingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNewRelicAlertMutingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func parseMutingRuleIDs(ids string) (*mutingRuleIDs, error) {
	split := strings.Split(ids, ":")

	accountID, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return nil, err
	}

	mutingRuleID, err := strconv.ParseInt(split[1], 10, 32)
	if err != nil {
		return nil, err
	}

	return &mutingRuleIDs{
		AccountID: int(accountID),
		ID:        int(mutingRuleID),
	}, nil
}

type mutingRuleIDs struct {
	AccountID int
	ID        int
}

func (w *mutingRuleIDs) String() string {
	return fmt.Sprintf("%d:%d", w.AccountID, w.ID)
}
