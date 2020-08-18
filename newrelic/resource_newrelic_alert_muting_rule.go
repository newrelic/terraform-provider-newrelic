package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
							Type:        schema.TypeList, //set or list
							Optional:    true,
							Description: "The individual MutingRuleConditions within the group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"accountId", "conditionId", "policyId", "policyName", "conditionName", "conditionType", "conditionRunbookUrl", "product", "targetId", "targetName", "nrqlEventType", "tag", "nrqlQuery"}, false),
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

	// conditions := alerts.MutingRuleCondition{}

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

	return condition
}

func expandMutingRuleConditionGroup(cfg map[string]interface{}) alerts.MutingRuleConditionGroup {

	conditionGroup := alerts.MutingRuleConditionGroup{}
	var expandedConditions []alerts.MutingRuleCondition

	log.Print("\n\n **************************** \n")
	log.Printf("\n CONDITIONS TYPE:  %T \n", cfg["conditions"])
	log.Print("\n **************************** \n\n")

	conditions := cfg["conditions"].([]interface{}) //need to fix

	for _, c := range conditions {
		// cast c to map string
		var y = expandMutingRuleCondition(c)
		expandedConditions = append(expandedConditions, y)
	}
	// sort.SliceStable(expandedConditions, func(i, j int) bool {
	// 	return expandedConditions[i].Attribute < expandedConditions[j].Attribute
	// })
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

	if e, ok := d.GetOk("condition"); ok { // log
		createInput.Condition = expandMutingRuleConditionGroup(e.([]interface{})[0].(map[string]interface{}))
	}

	return createInput
}

func resourceNewRelicAlertMutingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	createInput := expandMutingRuleCreateInput(d)

	fmt.Println("====================")
	fmt.Printf("createInput: %+v \n", createInput)
	fmt.Println("====================")

	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Creating New Relic MutingRule alerts")

	created, err := client.Alerts.CreateMutingRule(accountID, createInput)
	if err != nil {
		return err
	}

	// fmt.Println("====================")
	// fmt.Printf("CREATE RESULT: %+v \n", created)
	// fmt.Println("====================")

	d.SetId(serializeIDs([]int{accountID, created.ID}))

	return resourceNewRelicAlertMutingRuleRead(d, meta)
}

func flattenMutingRuleCondition(in alerts.MutingRuleCondition) map[string]interface{} {
	return map[string]interface{}{
		"attribute": in.Attribute,
		"operator":  in.Operator,
		"values":    in.Values,
	}
}

func flattenMutingRuleConditionGroup(in alerts.MutingRuleConditionGroup) []map[string]interface{} {
	condition := []map[string]interface{}{
		{
			"operator": in.Operator,
		},
	}

	for _, src := range in.Conditions {
		dst := map[string]interface{}{
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

	log.Print("\n\n **************************** \n")
	log.Printf("\n flattenMutingRule:  %+v \n", d.Get("condition"))
	log.Print("\n **************************** \n\n")

	return nil
}

func resourceNewRelicAlertMutingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic MutingRule alerts")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	accountID := ids[0]
	mutingRuleID := ids[1]

	mutingRule, err := client.Alerts.GetMutingRule(accountID, mutingRuleID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	// fmt.Println("====================")
	// fmt.Printf("READ RESULT: %+v \n", mutingRule)
	// fmt.Println("====================")

	return flattenMutingRule(mutingRule, d)

}

func expandMutingRuleUpdateInput(d *schema.ResourceData) alerts.MutingRuleUpdateInput {
	updateInput := alerts.MutingRuleUpdateInput{
		Enabled:     d.Get("enabled").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if e, ok := d.GetOk("condition"); ok { //detect if its nil
		fmt.Println("====================")
		fmt.Printf("e: %+v \n", e)
		fmt.Println("====================")
		x := expandMutingRuleConditionGroup(e.([]interface{})[0].(map[string]interface{}))

		fmt.Println("====================")
		fmt.Printf("x: %+v \n", x)
		fmt.Println("====================")
		updateInput.Condition = &x
	}

	return updateInput
}

func resourceNewRelicAlertMutingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	updateInput := expandMutingRuleUpdateInput(d)

	fmt.Println("====================")
	fmt.Printf("updateInput: %+v \n", updateInput)
	fmt.Println("====================")

	log.Printf("[INFO] Updating New Relic One alert muting rule.")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	accountID := ids[0]
	mutingRuleID := ids[1]

	_, err = client.Alerts.UpdateMutingRule(accountID, mutingRuleID, updateInput)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
	}
	// 	return err

	// d.SetId(ids.String())

	return resourceNewRelicAlertMutingRuleRead(d, meta)
}

func resourceNewRelicAlertMutingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One muting rule alert.")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	accountID := ids[0]
	mutingRuleID := ids[1]

	client.Alerts.DeleteMutingRule(accountID, mutingRuleID)

	return nil
}

// func flattenMutingRuleValues(in []alerts.MutingRuleCondition) []string {
// 	var values []map[string]interface{}
//
// 	for _, src := range *in {
// 		dst := map[string]interface{}{
// 			"value": src.Value,
// 		}
// 		values = append(values, dst)
// 	}
//
// 	return values
// }
