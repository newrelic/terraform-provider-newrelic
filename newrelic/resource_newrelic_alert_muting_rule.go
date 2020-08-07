package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func resourceNewRelicAlertMutingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	// condition, err := expandAlertCondition(d)

	accountID := d.Get("account_id").(int)

	log.Printf("[INFO] Creating New Relic MutingRule alerts")

	created, err = client.Alerts.CreateMutingRule(accountID, *created)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{accountID, created.ID}))

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
