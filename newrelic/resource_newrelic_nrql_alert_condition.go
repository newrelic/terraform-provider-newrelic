package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func mergeSchemas(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	res := map[string]*schema.Schema{}
	for _, m := range schemas {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

func resourceNewRelicNrqlAlertCondition() *schema.Resource {
	// Base schema (common attributes and blocks between old and new API)
	schemaBase := map[string]*schema.Schema{
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
			Description: "Whether or not to enable the alert condition.",
		},
		// The "outlier" type does NOT exist in NerdGraph yet (needs custom validation)
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "static",
			ValidateFunc: validation.StringInSlice([]string{"static", "outlier", "baseline"}, false),
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

					// Does not exist in NerdGraph. Handle this scenario when NerdGraph is called.
					"since_value": {
						Deprecated: "use `evaluation_offset` attribute instead",
						Type:       schema.TypeString,
						Optional:   true, // Used to be required
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
						// ConflictsWith: []string{"evaluation_offset"},
					},

					// New field in NerdGraph. Equivalent to `since_value`.
					"evaluation_offset": {
						Type:          schema.TypeInt,
						Optional:      true,
						Default:       1,
						ValidateFunc:  validation.IntBetween(1, 20),
						ConflictsWith: []string{"nrql.0.since_value"},
					},
				},
			},
		},
		"term": {
			Type:        schema.TypeSet,
			Description: "A list of terms for this condition. ",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Maps to `thresholdDuration` in NerdGraph and values are in seconds, not minutes.
					// Validation is different in NerdGraph - Value must be within 120-3600 seconds (2-60 minutes) and a multiple of 60 for BASELINE conditions.
					// Convert to seconds when using NerdGraph
					"duration": {
						Type:         schema.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 120),
						Description:  "In minutes, must be in the range of 1 to 120, inclusive.",
					},
					// Value must be uppercase when using NerdGraph
					"operator": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "equal",
						ValidateFunc: validation.StringInSlice([]string{"above", "below", "equal"}, false),
						Description:  "One of (above, below, equal). Defaults to equal.",
					},
					// Value must be uppercase when using NerdGraph
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

					// Does not exist in NerdGraph
					"time_function": {
						Deprecated:   "use `threshold_occurrences` attribute instead",
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
						Description:  "Valid values are: 'all' or 'any'",
					},

					// NerdGraph only. Seems to be similar to `time_function`
					"threshold_occurrences": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"ALL", "AT_LEAST_ONCE"}, false),
						Description:  "Valid values are: 'ALL' or 'AT_LEAST_ONCE'",
					},
				},
			},
			Required: true,
			MinItems: 1,
		},
	}

	// Old
	oldSchemaFields := map[string]*schema.Schema{
		// Outlier ONLY
		"expected_groups": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of expected groups when using outlier detection.",
		},
		// Outlier ONLY
		"ignore_overlap": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to look for a convergence of groups when using outlier detection.",
		},
		"violation_time_limit_seconds": {
			Deprecated:   "use `violation_time_limit` attribute instead",
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntInSlice([]int{3600, 7200, 14400, 28800, 43200, 86400}),
			Description:  "Sets a time limit, in seconds, that will automatically force-close a long-lasting violation after the time limit you select. Possible values are 3600, 7200, 14400, 28800, 43200, and 86400.",
		},
		// Exists in NerdGraph, but with different values. Figure out how to handle this.
		// Conflicts with `baseline_direction` when using NerdGraph
		"value_function": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "single_value",
			ValidateFunc: validation.StringInSlice([]string{"single_value", "sum"}, false),
			Description:  "Valid values are: 'single_value' or 'sum'",
		},
	}

	// New
	newSchemaFields := map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The New Relic account ID for managing your NRQL alert conditions.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the NRQL alert condition.",
		},
		"violation_time_limit": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"ONE_HOUR", "TWO_HOURS", "FOUR_HOURS", "EIGHT_HOURS", "TWELVE_HOURS", "TWENTY_FOUR_HOURS"}, false),
			Description:  "Sets a time limit, in hours, that will automatically force-close a long-lasting violation after the time limit you select. Possible values are ONE_HOUR, TWO_HOURS, FOUR_HOURS, EIGHT_HOURS, TWELVE_HOURS, TWENTY_FOUR_HOURS.",
		},
		// Baseline ONLY
		"baseline_direction": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "The baseline direction of a baseline NRQL alert condition. Valid values are: 'LOWER_ONLY', 'UPPER_AND_LOWER', 'UPPER_ONLY'",
			ValidateFunc:  validation.StringInSlice([]string{"LOWER_ONLY", "UPPER_AND_LOWER", "UPPER_ONLY"}, false),
			ConflictsWith: []string{"value_function"},
		},
	}

	nrqlAlertConditionSchema := mergeSchemas(schemaBase, oldSchemaFields, newSchemaFields)

	return &schema.Resource{
		Create: resourceNewRelicNrqlAlertConditionCreate,
		Read:   resourceNewRelicNrqlAlertConditionRead,
		Update: resourceNewRelicNrqlAlertConditionUpdate,
		Delete: resourceNewRelicNrqlAlertConditionDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: nrqlAlertConditionSchema,
	}
}

// Selects the proper accountID for usage within a resource. Account IDs provided
// within a `resource` block will override a `provider` block account ID. This ensures
// resources can be scoped to specific accounts. Bear in mind those accounts must be
// accessible with the provided Personal API Key (APIKS).
func selectAccountID(providerCondig *ProviderConfig, d *schema.ResourceData) int {
	resourceAccountID := d.Get("account_id").(int)

	if resourceAccountID != 0 {
		return resourceAccountID
	}

	return providerCondig.AccountID
}

func resourceNewRelicNrqlAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	fmt.Println(" ")
	fmt.Println("****************************")
	fmt.Printf("\n\n CONFIG   Acct ID: %+v \n", accountID)
	fmt.Printf(" CONFIG   API Key: %+v \n\n", providerConfig.PersonalAPIKey)
	fmt.Println("****************************")
	fmt.Println(" ")

	if accountID != 0 && providerConfig.PersonalAPIKey != "" {
		conditionInput, err := expandNrqlAlertConditionInput(d)
		if err != nil {
			return err
		}

		log.Printf("[INFO] Creating New Relic NRQL alert condition %s via NerdGraph API", conditionInput.Name)

		fmt.Println(" ")
		fmt.Println("****************************")
		fmt.Printf("\n\n conditionInput: %+v \n\n", conditionInput)
		fmt.Printf("\n\n Is Baseline?    %+v \n\n", conditionInput.Type == "baseline")
		fmt.Println("****************************")
		fmt.Println(" ")

		policyID := d.Get("policy_id").(int)

		if conditionInput.Type == "baseline" {
			_, err = client.Alerts.CreateNrqlConditionBaselineMutation(accountID, policyID, *conditionInput)
			if err != nil {
				return err
			}
		}

		// if conditionInput.Type == "static" {
		// 	_, err = client.Alerts.CreateNrqlConditionStaticMutation(accountID, policyID, *conditionInput)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}

	// condition := expandNrqlAlertConditionStruct(d)
	// policyID := d.Get("policy_id").(int)

	// log.Printf("[INFO] Creating New Relic NRQL alert condition %s", condition.Name)

	// condition, err := client.Alerts.CreateNrqlCondition(policyID, *condition)
	// if err != nil {
	// 	return err
	// }

	// d.SetId(serializeIDs([]int{policyID, condition.ID}))

	// return resourceNewRelicNrqlAlertConditionRead(d, meta)
	return errors.NewNotFound("Testing")
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
