package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNewRelicAlertMutingRule() *schema.Resource {
	validAlertConditionTypes := make([]string, 0, len(alertConditionTypes))
	for k := range alertConditionTypes {
		validAlertConditionTypes = append(validAlertConditionTypes, k)
	}

	return &schema.Resource{
		Create: resourceNewRelicAlertMutingRuleCreate,
		Read:   resourceNewRelicAlertMutingRuleRead,
		Update: resourceNewRelicAlertMutingRuleUpdate,
		Delete: resourceNewRelicAlertMutingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Input settings for the muting rule.",
			},
			"accountId": {
				Type:        schema.TypeAlertsMutingRuleInput, // Not sure about this..
				Optional:    true,
				Description: "The muting rule's account Id.",
			},
			"condition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The condition that defines which violations to target.",
			},
			"createdAt": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier for the muting rule",
			},
			"createdBy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier for the muting rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the muting rule.",
			},

			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier for the muting rule.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the MutingRule.",
			},
			// "schedule": {
			// 	Type: schema.TypeSet,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"startTile": {
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 				Description: "The unique identifier for the muting rule",
			// 			},
			// },

			"updatedAt": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier for the muting rule",
			},
			"updatedBy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier for the muting rule",
			},
			// "enabled": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Description: "The unique identifier for the muting rule",
			// },

		},
	}
}

func resourceNewRelicAlertMutingRuleCreate(d *schema.ResourceData, meta interface{}) error {}

func resourceNewRelicAlertMutingRuleRead(d *schema.ResourceData, meta interface{}) error {}

func resourceNewRelicAlertMutingRuleUpdate(d *schema.ResourceData, meta interface{}) error {}

func resourceNewRelicAlertMutingRuleDelete(d *schema.ResourceData, meta interface{}) error {}
