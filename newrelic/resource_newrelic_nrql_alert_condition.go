package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// termSchema returns the schema used for a critical or warning term priority.
func termSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Maps to `thresholdDuration` in NerdGraph and values are in seconds, not minutes.
			// Validation is different in NerdGraph - Value must be within 120-3600 seconds (2-60 minutes) and a multiple of 60 for BASELINE conditions.
			// Convert to seconds when using NerdGraph
			"duration": {
				Deprecated:   "use `threshold_duration` attribute instead",
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "In minutes, must be in the range of 1 to 120 (inclusive).",
				ValidateFunc: validation.IntBetween(1, 120),
			},
			// Value must be uppercase when using NerdGraph
			"operator": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "equals",
				Description:  "One of (above, below, equals). Defaults to 'equals'.",
				ValidateFunc: validation.StringInSlice([]string{"above", "below", "equals"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"threshold": {
				Type:         schema.TypeFloat,
				Required:     true,
				Description:  "Must be 0 or greater. For baseline conditions must be in range [1, 1000].",
				ValidateFunc: float64Gte(0.0),
			},
			// Does not exist in NerdGraph. Equivalent to `threshold_occurrences`,
			// but with different wording.
			"time_function": {
				Deprecated:   "use `threshold_occurrences` attribute instead",
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Valid values are: 'all' or 'any'",
				ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
			},
			// NerdGraph only. Equivalent to `time_function`,
			// but with slightly different wording.
			// i.e. `any` (old) vs `AT_LEAST_ONCE` (new)
			"threshold_occurrences": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The criteria for how many data points must be in violation for the specified threshold duration. Valid values are: 'ALL' or 'AT_LEAST_ONCE' (case insensitive).",
				ValidateFunc: validation.StringInSlice([]string{"ALL", "AT_LEAST_ONCE"}, true),
				StateFunc: func(v interface{}) string {
					// Always store lowercase to prevent state drift
					return strings.ToLower(v.(string))
				},
			},
			// NerdGraph only. Equivalent to `duration`, but in seconds
			"threshold_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The duration, in seconds, that the threshold must violate in order to create a violation. Value must be a multiple of the 'aggregation_window' (which has a default of 60 seconds). Value must be within 120-3600 seconds for baseline and outlier conditions, within 120-7200 seconds for static conditions with the sum value function, and within 60-7200 seconds for static conditions with the single_value value function.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					minVal := 60
					maxVal := 7200

					// Value must be a factor of 60.
					if v%60 != 0 {
						errs = append(errs, fmt.Errorf("%q must be a factor of %d, got: %d", key, minVal, v))
					}

					// This validation is a top-level validation check.
					// Static conditions with a single_value value function must be within range [60, 7200].
					// Static conditions with a sum value function must be within range [120, 7200].
					// Baseline conditions must be within range [120, 3600].
					// Outlier conditions must be within range [120, 3600].
					if v < minVal || v > maxVal {
						errs = append(errs, fmt.Errorf("%q must be between %d and %d inclusive, got: %d", key, minVal, maxVal, v))
					}

					return
				},
			},
		},
	}
}

func termSchemaDeprecated() *schema.Resource {
	rec := termSchema()

	prioritySchema := &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "critical",
		Description:  "One of (critical, warning). Defaults to 'critical'. At least one condition term must have priority set to 'critical'.",
		ValidateFunc: validation.StringInSlice([]string{"critical", "warning"}, false),
	}

	rec.Schema["priority"] = prioritySchema

	return rec
}

func resourceNewRelicNrqlAlertCondition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNrqlAlertConditionCreate,
		ReadContext:   resourceNewRelicNrqlAlertConditionRead,
		UpdateContext: resourceNewRelicNrqlAlertConditionUpdate,
		DeleteContext: resourceNewRelicNrqlAlertConditionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateWithMetadata(2, "type"),
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
			// Note: The "outlier" type does NOT exist in NerdGraph yet
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "static",
				Description: "The type of NRQL alert condition to create. Valid values are: 'static', 'baseline', 'outlier' (deprecated).",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					valueString := val.(string)

					v := validation.StringInSlice([]string{"static", "outlier", "baseline"}, false)
					warns, errs = v(valueString, key)

					if valueString == "outlier" {
						warns = append(warns, "We're removing outlier conditions Feb 1, 2022. More Info: https://docs.newrelic.com/docs/alerts-applied-intelligence/transition-guide/#outlier")
					}
					return
				},
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
						"since_value": {
							Deprecated:    "use `signal.aggregation_method` attribute instead",
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "NRQL queries are evaluated in one-minute time windows. The start time depends on the value you provide in the NRQL condition's `since_value`.",
							ConflictsWith: []string{"nrql.0.evaluation_offset"},
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
						},
						// Equivalent to `since_value`.
						"evaluation_offset": {
							Deprecated:    "use `signal.aggregation_method` attribute instead",
							Type:          schema.TypeInt,
							Optional:      true,
							Description:   "NRQL queries are evaluated in one-minute time windows. The start time depends on the value you provide in the NRQL condition's `evaluation_offset`.",
							ConflictsWith: []string{"nrql.0.since_value"},
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 20 {
									errs = append(errs, fmt.Errorf("%q must be between 0 and 20 inclusive, got: %d", key, v))
								}
								return
							},
						},
					},
				},
			},
			"term": {
				Type:          schema.TypeSet,
				MinItems:      1,
				MaxItems:      2,
				Optional:      true,
				Description:   "A set of terms for this condition. Max 2 terms allowed - at least one 1 critical term and 1 optional warning term.",
				Elem:          termSchemaDeprecated(),
				ConflictsWith: []string{"critical", "warning"},
				Deprecated:    "use `critical` and `warning` attributes instead",
			},
			"critical": {
				Type:          schema.TypeList,
				MinItems:      1,
				MaxItems:      1,
				Optional:      true,
				Elem:          termSchema(),
				Description:   "A condition term with priority set to critical.",
				ConflictsWith: []string{"term"},
			},
			"warning": {
				Type:          schema.TypeList,
				MinItems:      1,
				MaxItems:      1,
				Optional:      true,
				Elem:          termSchema(),
				Description:   "A condition term with priority set to warning.",
				ConflictsWith: []string{"term"},
			},
			// Outlier ONLY
			"expected_groups": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of expected groups when using outlier detection.",
			},
			// Outlier ONLY
			"ignore_overlap": {
				Deprecated:    "use `open_violation_on_group_overlap` attribute instead, but use the inverse of your boolean - e.g. if ignore_overlap = false, use open_violation_on_group_overlap = true",
				Type:          schema.TypeBool,
				Optional:      true,
				Description:   "Whether to look for a convergence of groups when using outlier detection.",
				ConflictsWith: []string{"open_violation_on_group_overlap"},
			},
			"violation_time_limit_seconds": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "Sets a time limit, in seconds, that will automatically force-close a long-lasting violation after the time limit you select.  Must be in the range of 300 to 2592000 (inclusive)",
				ConflictsWith: []string{"violation_time_limit"},
				ValidateFunc:  validation.IntBetween(300, 2592000),
			},
			// Exists in NerdGraph, but with different values. Conversion
			// between new:old and old:new is handled via maps in structures file.
			// Conflicts with `baseline_direction` when using NerdGraph.
			"value_function": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Valid values are: 'single_value' or 'sum'",
				ValidateFunc: validation.StringInSlice([]string{"single_value", "sum"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // Case fold this attribute when diffing
				},
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID for managing your NRQL alert conditions.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the NRQL alert condition.",
			},
			"violation_time_limit": {
				Type:          schema.TypeString,
				Deprecated:    "use `violation_time_limit_seconds` attribute instead",
				Optional:      true,
				Computed:      true,
				Description:   "Sets a time limit, in hours, that will automatically force-close a long-lasting violation after the time limit you select. Possible values are 'ONE_HOUR', 'TWO_HOURS', 'FOUR_HOURS', 'EIGHT_HOURS', 'TWELVE_HOURS', 'TWENTY_FOUR_HOURS', 'THIRTY_DAYS' (case insensitive).",
				ConflictsWith: []string{"violation_time_limit_seconds"},
				ValidateFunc:  validation.StringInSlice([]string{"ONE_HOUR", "TWO_HOURS", "FOUR_HOURS", "EIGHT_HOURS", "TWELVE_HOURS", "TWENTY_FOUR_HOURS", "THIRTY_DAYS"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // Case fold this attribute when diffing
				},
			},
			"open_violation_on_expiration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to create a new violation to capture that the signal expired.",
			},
			"close_violations_on_expiration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to close all open violations when the signal expires.",
			},
			"aggregation_window": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The duration of the time window used to evaluate the NRQL query, in seconds.",
			},
			"slide_by": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The duration of overlapping timewindows used to smooth the chart line, in seconds. Must be a factor of `aggregation_window` and less than the aggregation window. It should be greater or equal to 30 seconds if `aggregation_window` is less than or equal to 3600 seconds, or greater or eqaul to `aggregation_window / 120` if `aggregation_window` is greater than 3600 seconds.",
			},
			"expiration_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The amount of time (in seconds) to wait before considering the signal expired.",
			},
			"fill_option": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Which strategy to use when filling gaps in the signal. If static, the 'fill value' will be used for filling gaps in the signal. Valid values are: 'NONE', 'LAST_VALUE', or 'STATIC' (case insensitive).",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "LAST_VALUE", "STATIC"}, true),
				StateFunc: func(v interface{}) string {
					// Always store lowercase to prevent state drift
					return strings.ToLower(v.(string))
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Assume that empty string and 'none' are the same for diff purposes due to API defaults
					return (old == "" || old == "none") && (new == "" || new == "none")
				},
			},
			"fill_value": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Description:  "If using the 'static' fill option, this value will be used for filling gaps in the signal.",
				RequiredWith: []string{"fill_option"},
			},
			"aggregation_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CADENCE", "EVENT_FLOW", "EVENT_TIMER"}, true),
				Description:  "The method that determines when we consider an aggregation window to be complete so that we can evaluate the signal for violations. Default is CADENCE.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// If a value is not provided and the condition uses the default value, don't show a diff
					return (old == "event_flow" && new == "") || strings.EqualFold(old, new)
				},
			},
			"aggregation_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "How long we wait for data that belongs in each aggregation window. Depending on your data, a longer delay may increase accuracy but delay notifications. Use aggregationDelay with the EVENT_FLOW and CADENCE aggregation methods.",
				RequiredWith: []string{"aggregation_method"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// If a value is not provided and the condition uses the default value, don't show a diff
					oldInt, _ := strconv.ParseInt(old, 0, 8)
					newInt, _ := strconv.ParseInt(new, 0, 8)
					return oldInt == 120 && (newInt == 0)
				},
			},
			"aggregation_timer": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "How long we wait after each data point arrives to make sure we've processed the whole batch. Use aggregationTimer with the EVENT_TIMER aggregation method.",
				RequiredWith: []string{"aggregation_method"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// If a value is not provided and the condition uses the default value, don't show a diff
					oldInt, _ := strconv.ParseInt(old, 0, 8)
					newInt, _ := strconv.ParseInt(new, 0, 8)
					return oldInt == 120 && (newInt == 0)
				},
			},
			// Baseline ONLY
			"baseline_direction": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The baseline direction of a baseline NRQL alert condition. Valid values are: 'LOWER_ONLY', 'UPPER_AND_LOWER', 'UPPER_ONLY' (case insensitive).",
				ConflictsWith: []string{"value_function"},
				ValidateFunc:  validation.StringInSlice([]string{"LOWER_ONLY", "UPPER_AND_LOWER", "UPPER_ONLY"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // Case fold this attribute when diffing
				},
			},
			// Outlier ONLY
			"open_violation_on_group_overlap": {
				Type:          schema.TypeBool,
				Optional:      true,
				Description:   "Whether overlapping groups should produce a violation.",
				ConflictsWith: []string{"ignore_overlap"},
			},
		},
	}
}

func resourceNewRelicNrqlAlertConditionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)
	policyID := strconv.Itoa(d.Get("policy_id").(int))

	conditionInput, err := expandNrqlAlertConditionCreateInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic NRQL alert condition %s via NerdGraph API", conditionInput.Name)

	var condition *alerts.NrqlAlertCondition

	switch d.Get("type").(string) {
	case "baseline":
		condition, err = client.Alerts.CreateNrqlConditionBaselineMutationWithContext(ctx, accountID, policyID, *conditionInput)
	case "static":
		condition, err = client.Alerts.CreateNrqlConditionStaticMutationWithContext(ctx, accountID, policyID, *conditionInput)
	case "outlier":
		condition, err = client.Alerts.CreateNrqlConditionOutlierMutationWithContext(ctx, accountID, policyID, *conditionInput)
	}

	var diags diag.Diagnostics

	if graphQLError, ok := err.(*alerts.GraphQLErrorResponse); ok {
		for _, e := range graphQLError.Errors {
			var message string = e.Message
			var errorClass string = e.Extensions.ErrorClass
			var validationErrors = e.Extensions.ValidationErrors

			if len(validationErrors) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  message + ": " + errorClass,
				})
			} else {
				for _, validationError := range validationErrors {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  message + ": " + errorClass,
						Detail:   validationError.Name + ": " + validationError.Reason,
					})
				}
			}
		}

		return diags
	}

	conditionID, err := strconv.Atoi(condition.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serializeIDs([]int{d.Get("policy_id").(int), conditionID}))

	return resourceNewRelicNrqlAlertConditionRead(ctx, d, meta)
}

func resourceNewRelicNrqlAlertConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic NRQL alert condition %s", d.Id())

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := ids[0]
	conditionID := ids[1]

	_, err = client.Alerts.QueryPolicyWithContext(ctx, accountID, strconv.Itoa(policyID))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	nrqlCondition, err := client.Alerts.GetNrqlConditionQueryWithContext(ctx, accountID, strconv.Itoa(conditionID))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenNrqlAlertCondition(accountID, nrqlCondition, d))
}

func resourceNewRelicNrqlAlertConditionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	conditionID := strconv.Itoa(ids[1])

	conditionInput, err := expandNrqlAlertConditionUpdateInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	switch d.Get("type").(string) {
	case "baseline":
		_, err = client.Alerts.UpdateNrqlConditionBaselineMutationWithContext(ctx, accountID, conditionID, *conditionInput)
	case "static":
		_, err = client.Alerts.UpdateNrqlConditionStaticMutationWithContext(ctx, accountID, conditionID, *conditionInput)
	case "outlier":
		_, err = client.Alerts.UpdateNrqlConditionOutlierMutationWithContext(ctx, accountID, conditionID, *conditionInput)
	}

	var diags diag.Diagnostics

	if graphQLError, ok := err.(*alerts.GraphQLErrorResponse); ok {
		for _, e := range graphQLError.Errors {
			var message string = e.Message
			var errorClass string = e.Extensions.ErrorClass
			var validationErrors = e.Extensions.ValidationErrors

			if len(validationErrors) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  message + ": " + errorClass,
				})
			} else {
				for _, validationError := range validationErrors {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  message + ": " + errorClass,
						Detail:   validationError.Name + ": " + validationError.Reason,
					})
				}
			}
		}

		return diags
	}

	return resourceNewRelicNrqlAlertConditionRead(ctx, d, meta)
}

func resourceNewRelicNrqlAlertConditionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	conditionID := strconv.Itoa(ids[1])

	log.Printf("[INFO] Deleting New Relic NRQL alert condition %v", conditionID)

	_, err = client.Alerts.DeleteNrqlConditionMutationWithContext(ctx, accountID, conditionID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
